package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoError "github.com/ladmakhi81/learnup/internals/video/error"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/db/repositories"
	"github.com/ladmakhi81/learnup/shared/types"
	"github.com/ladmakhi81/learnup/shared/utils"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type VideoService interface {
	UpdateURLAndDuration(dto dtoreq.UpdateURLAndDurationVideoReqDto) (*entities.Video, error)
	CreateCompleteUploadVideoNotification(videoID uint) error
	Encode(ctx context.Context, dto dtoreq.EncodeVideoReqDto) (string, error)
	CalculateDuration(ctx context.Context, dto dtoreq.CalculateVideoDurationReqDto) (string, error)
	Verify(admin *entities.User, videoId uint) error
	FindVideosByCourseID(courseID uint) ([]*entities.Video, error)
}

type videoService struct {
	unitOfWork  db.UnitOfWork
	minioClient contracts.Storage
	ffmpegSvc   contracts.Ffmpeg
	logSvc      contracts.Log
}

func NewVideoSvc(
	unitOfWork db.UnitOfWork,
	minioClient contracts.Storage,
	ffmpegSvc contracts.Ffmpeg,
	logSvc contracts.Log,
) VideoService {
	return &videoService{
		unitOfWork:  unitOfWork,
		minioClient: minioClient,
		ffmpegSvc:   ffmpegSvc,
		logSvc:      logSvc,
	}
}

func (svc videoService) CreateCompleteUploadVideoNotification(videoID uint) error {
	const operationName = "videoService.CreateCompleteUploadVideoNotification"
	course, err := svc.unitOfWork.CourseRepo().GetByVideoID(videoID)
	if err != nil {
		return types.NewServerError("Error in fetching course by video id", operationName, err)
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
	if err := svc.unitOfWork.NotificationRepo().Create(notification); err != nil {
		return types.NewServerError("Error in creating notification", operationName, err)
	}
	return nil
}

func (svc videoService) UpdateURLAndDuration(dto dtoreq.UpdateURLAndDurationVideoReqDto) (*entities.Video, error) {
	const operationName = "videoService.UpdateURLAndDuration"
	video, err := svc.unitOfWork.VideoRepo().GetByID(dto.ID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching video by id", operationName, err)
	}
	if video == nil {
		return nil, videoError.Video_NotFound
	}
	video.URL = dto.URL
	video.Duration = &dto.Duration
	video.Status = entities.VideoStatus_Done
	if err := svc.unitOfWork.VideoRepo().Update(video); err != nil {
		return nil, types.NewServerError("Error in updating the video", operationName, err)
	}
	return video, nil
}

func (svc videoService) CalculateDuration(ctx context.Context, dto dtoreq.CalculateVideoDurationReqDto) (string, error) {
	const operationName = "videoService.CalculateDuration"
	file, err := svc.minioClient.GetFileReader(ctx, "videos", dto.ObjectId)
	if err != nil {
		return "", types.NewServerError("Error in get file from minio", operationName, err)
	}
	durationStr, err := svc.ffmpegSvc.GetVideoDuration(file)
	if err != nil {
		return "", types.NewServerError("Error in calculating video duration", operationName, err)
	}
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return "", types.NewServerError("Error in converting duration into number", operationName, err)
	}
	hours := int(duration) / 3600
	minutes := (int(duration) % 3600) / 60
	seconds := int(math.Mod(duration, 60))
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds), nil
}

func (svc videoService) Encode(ctx context.Context, dto dtoreq.EncodeVideoReqDto) (string, error) {
	const operationName = "videoService.Encode"
	log.Println("encode function execute")
	// encode
	file, err := svc.minioClient.GetFileReader(ctx, "videos", dto.ObjectId)
	if err != nil {
		return "", types.NewServerError("Error in get file from minio", operationName, err)
	}
	storeLocation, err := svc.ffmpegSvc.EncodeVideo(file)
	if err != nil {
		return "", types.NewServerError("Error in encoding video file", operationName, err)
	}

	// move from local to minio
	dirFiles, err := os.ReadDir(storeLocation)
	if err != nil {
		return "", types.NewServerError("Error in finding directories of encoded files", operationName, err)
	}
	contentTypes := map[string]string{
		".ts":   "video/mp2t",
		".m3u8": "application/vnd.apple.mpegurl",
	}
	encodedFilePath := uuid.NewString()
	for _, dirFile := range dirFiles {
		file, err := os.ReadFile(path.Join(storeLocation, dirFile.Name()))
		if err != nil {
			return "", types.NewServerError("Error in finding file of directories of encoded files", operationName, err)
		}
		fileExt := filepath.Ext(dirFile.Name())
		currentContentType := contentTypes[fileExt]
		if _, err := svc.minioClient.UploadFileByContent(
			ctx,
			"videos",
			fmt.Sprintf("%s/%s", encodedFilePath, dirFile.Name()),
			currentContentType,
			file,
		); err != nil {
			return "", types.NewServerError("Error in storing encoded video into storage", operationName, err)
		}
	}

	// remove unused files
	if err := os.RemoveAll(storeLocation); err != nil {
		return "", types.NewServerError("Error in deleting file from disk", operationName, err)
	}
	if err := svc.minioClient.DeleteObject(context.TODO(), "videos", dto.ObjectId); err != nil {
		return "", types.NewServerError("Error in deleting file from minio", operationName, err)
	}
	if err := svc.minioClient.DeleteObject(context.TODO(), "videos", fmt.Sprintf("%s.info", dto.ObjectId)); err != nil {
		return "", types.NewServerError("Error in deleting file from minio", operationName, err)
	}
	return encodedFilePath, nil
}

func (svc videoService) Verify(admin *entities.User, videoId uint) error {
	const operationName = "videoService.Verify"
	video, err := svc.unitOfWork.VideoRepo().GetByID(videoId, nil)
	if err != nil {
		return types.NewServerError("Error in fetching video by id", operationName, err)
	}
	if video == nil {
		return videoError.Video_NotFound
	}
	if video.IsVerified {
		return nil
	}
	video.IsVerified = true
	video.VerifiedDate = utils.Now()
	video.VerifiedById = &admin.ID
	if err := svc.unitOfWork.VideoRepo().Update(video); err != nil {
		return types.NewServerError("Error in verify video", operationName, err)
	}
	return nil
}

func (svc videoService) FindVideosByCourseID(courseID uint) ([]*entities.Video, error) {
	const operationName = "videoService.FindVideosByCourseID"
	course, err := svc.unitOfWork.CourseRepo().GetByID(courseID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching course by id", operationName, err)
	}
	if course == nil {
		return nil, courseError.Course_NotFound
	}
	videos, err := svc.unitOfWork.VideoRepo().GetAll(
		repositories.GetAllOptions{
			Conditions: map[string]any{
				"course_id": courseID,
			},
			Relations: []string{"VerifiedBy"},
		},
	)
	if err != nil {
		return nil, types.NewServerError("Finding videos by course id throw error", operationName, err)
	}
	return videos, nil
}
