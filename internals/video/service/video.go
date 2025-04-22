package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type VideoService interface {
	UpdateURLAndDuration(dto dtoreq.UpdateURLAndDurationVideoReq) (*entities.Video, error)
	CreateCompleteUploadVideoNotification(videoID uint) error
	Encode(ctx context.Context, dto dtoreq.EncodeVideoReq) (string, error)
	CalculateDuration(ctx context.Context, dto dtoreq.CalculateVideoDurationReq) (string, error)
	Verify(authContext any, videoId uint) error
	FindVideosByCourseID(courseID uint) ([]*entities.Video, error)
}

type VideoServiceImpl struct {
	repo           *db.Repositories
	minioClient    contracts.Storage
	ffmpegSvc      contracts.Ffmpeg
	logSvc         contracts.Log
	translationSvc contracts.Translator
}

func NewVideoServiceImpl(
	repo *db.Repositories,
	minioClient contracts.Storage,
	ffmpegSvc contracts.Ffmpeg,
	logSvc contracts.Log,
	translationSvc contracts.Translator,
) *VideoServiceImpl {
	return &VideoServiceImpl{
		repo:           repo,
		minioClient:    minioClient,
		ffmpegSvc:      ffmpegSvc,
		logSvc:         logSvc,
		translationSvc: translationSvc,
	}
}

func (svc VideoServiceImpl) CreateCompleteUploadVideoNotification(videoID uint) error {
	course, courseErr := svc.repo.CourseRepo.GetByVideoID(videoID)
	if courseErr != nil {
		return types.NewServerError(
			"Error in fetching course by video id",
			"VideoServiceImpl.CreateCompleteUploadVideoNotification",
			courseErr,
		)
	}
	notification := &entities.Notification{
		Type: entities.NotificationType_CompleteVideoUpload,
		Metadata: map[string]any{
			"videoId":     videoID,
			"courseId":    course.ID,
			"courseTitle": course.Name,
		},
		IsSeen: false,
		UserID: course.TeacherID,
	}
	if err := svc.repo.NotificationRepo.Create(notification); err != nil {
		return types.NewServerError(
			"Error in creating notification",
			"VideoServiceImpl.CreateCompleteUploadVideoNotification",
			err,
		)
	}
	return nil
}

func (svc VideoServiceImpl) UpdateURLAndDuration(dto dtoreq.UpdateURLAndDurationVideoReq) (*entities.Video, error) {
	video, videoErr := svc.repo.VideoRepo.GetByID(dto.ID, nil)
	if videoErr != nil {
		return nil, types.NewServerError(
			"Error in fetching video by id",
			"VideoServiceImpl.UpdateURLAndDuration",
			videoErr,
		)
	}
	if video == nil {
		return nil, types.NewNotFoundError(svc.translationSvc.Translate("video.errors.not_found"))
	}
	video.URL = dto.URL
	video.Duration = &dto.Duration
	video.Status = entities.VideoStatus_Done
	if err := svc.repo.VideoRepo.Update(video); err != nil {
		return nil, types.NewServerError(
			"Error in updating the video",
			"VideoServiceImpl.UpdateVideoURL",
			err,
		)
	}
	return video, nil
}

func (svc VideoServiceImpl) CalculateDuration(ctx context.Context, dto dtoreq.CalculateVideoDurationReq) (string, error) {
	file, fileErr := svc.minioClient.GetFileReader(ctx, "videos", dto.ObjectId)
	if fileErr != nil {
		return "", types.NewServerError(
			"Error in get file from minio",
			"VideoServiceImpl.CalculateDuration",
			fileErr,
		)
	}
	durationStr, durationErr := svc.ffmpegSvc.GetVideoDuration(file)
	if durationErr != nil {
		return "", types.NewServerError(
			"Error in calculating video duration",
			"VideoServiceImpl.CalculateDuration",
			durationErr,
		)
	}
	duration, durationErr := strconv.ParseFloat(durationStr, 64)
	if durationErr != nil {
		return "", types.NewServerError(
			"Error in converting duration into number",
			"VideoServiceImpl.CalculateDuration",
			durationErr,
		)
	}
	hours := int(duration) / 3600
	minutes := (int(duration) % 3600) / 60
	seconds := int(math.Mod(duration, 60))
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds), nil
}

func (svc VideoServiceImpl) Encode(ctx context.Context, dto dtoreq.EncodeVideoReq) (string, error) {
	log.Println("encode function execute")
	// encode
	file, fileErr := svc.minioClient.GetFileReader(ctx, "videos", dto.ObjectId)
	if fileErr != nil {
		return "", types.NewServerError(
			"Error in get file from minio",
			"VideoServiceImpl.Encode",
			fileErr,
		)
	}
	storeLocation, storeLocationErr := svc.ffmpegSvc.EncodeVideo(file)
	if storeLocationErr != nil {
		return "", types.NewServerError(
			"Error in encoding video file",
			"VideoServiceImpl.Encode",
			storeLocationErr,
		)
	}

	// move from local to minio
	dirFiles, dirErr := os.ReadDir(storeLocation)
	if dirErr != nil {
		return "", types.NewServerError(
			"Error in finding directories of encoded files",
			"VideoServiceImpl.Encode",
			dirErr,
		)
	}
	contentTypes := map[string]string{
		".ts":   "video/mp2t",
		".m3u8": "application/vnd.apple.mpegurl",
	}
	encodedFilePath := uuid.NewString()
	for _, dirFile := range dirFiles {
		file, fileErr := os.ReadFile(path.Join(storeLocation, dirFile.Name()))
		if fileErr != nil {
			return "", types.NewServerError(
				"Error in finding file of directories of encoded files",
				"VideoServiceImpl.Encode",
				fileErr,
			)
		}
		fileExt := filepath.Ext(dirFile.Name())
		currentContentType := contentTypes[fileExt]
		if _, err := svc.minioClient.UploadFileByContent(
			context.TODO(),
			"videos",
			fmt.Sprintf("%s/%s", encodedFilePath, dirFile.Name()),
			currentContentType,
			file,
		); err != nil {
			return "", types.NewServerError(
				"Error in storing encoded video into storage",
				"VideoServiceImpl.Encode",
				err,
			)
		}
	}

	// remove unused files
	if err := os.RemoveAll(storeLocation); err != nil {
		return "", types.NewServerError(
			"Error in deleting file from disk",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			err,
		)
	}
	if err := svc.minioClient.DeleteObject(context.TODO(), "videos", dto.ObjectId); err != nil {
		return "", types.NewServerError(
			"Error in deleting file from minio",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			err,
		)
	}
	if err := svc.minioClient.DeleteObject(context.TODO(), "videos", fmt.Sprintf("%s.info", dto.ObjectId)); err != nil {
		return "", types.NewServerError(
			"Error in deleting file from minio",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			err,
		)
	}
	return encodedFilePath, nil
}

func (svc VideoServiceImpl) Verify(authContext any, videoId uint) error {
	video, videoErr := svc.repo.VideoRepo.GetByID(videoId, nil)
	if videoErr != nil {
		return types.NewServerError(
			"Error in fetching video by id",
			"VideoServiceImpl.Verify",
			videoErr,
		)
	}
	if video == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("video.errors.not_found"),
		)
	}
	if video.IsVerified {
		return nil
	}
	adminClaim := authContext.(*types.TokenClaim)
	admin, adminErr := svc.repo.UserRepo.GetByID(adminClaim.UserID, nil)
	if adminErr != nil {
		return types.NewServerError(
			"Error in fetching logged in user",
			"VideoServiceImpl.Verify",
			adminErr,
		)
	}
	if admin == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.admin_not_found"),
		)
	}
	now := time.Now()
	video.IsVerified = true
	video.VerifiedDate = &now
	video.VerifiedById = &admin.ID
	if err := svc.repo.VideoRepo.Update(video); err != nil {
		return types.NewServerError(
			"Error in verify video",
			"VideoServiceImpl.Verify",
			err,
		)
	}
	return nil
}

func (svc VideoServiceImpl) FindVideosByCourseID(courseID uint) ([]*entities.Video, error) {
	course, courseErr := svc.repo.CourseRepo.GetByID(courseID, nil)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in fetching course by id",
			"VideoServiceImpl.FindVideosByCourseID",
			courseErr,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(svc.translationSvc.Translate("course.errors.not_found"))
	}
	videos, videosErr := svc.repo.VideoRepo.GetAll(
		repositories.GetAllOptions{
			Conditions: map[string]any{
				"course_id": courseID,
			},
			Relations: []string{"VerifiedBy"},
		},
	)
	if videosErr != nil {
		return nil, types.NewServerError(
			"Finding videos by course id throw error",
			"VideoServiceImpl.FetchByCourseId",
			videosErr,
		)
	}
	return videos, nil
}
