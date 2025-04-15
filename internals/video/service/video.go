package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	notificationReqDto "github.com/ladmakhi81/learnup/internals/notification/dto/req"
	notificationEntity "github.com/ladmakhi81/learnup/internals/notification/entity"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoEntity "github.com/ladmakhi81/learnup/internals/video/entity"
	"github.com/ladmakhi81/learnup/internals/video/repo"
	"github.com/ladmakhi81/learnup/pkg/ffmpeg"
	"github.com/ladmakhi81/learnup/pkg/logger"
	"github.com/ladmakhi81/learnup/pkg/storage"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/types"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type VideoService interface {
	AddVideo(dto *dtoreq.AddVideoToCourse) (*videoEntity.Video, error)
	FindByTitle(title string) (*videoEntity.Video, error)
	IsVideoTitleExist(title string) (bool, error)
	FindVideosByCourseID(courseID uint) ([]*videoEntity.Video, error)
	UpdateURLAndDuration(dto dtoreq.UpdateURLAndDurationVideoReq) (*videoEntity.Video, error)
	FindById(id uint) (*videoEntity.Video, error)
	CreateCompleteUploadVideoNotification(videoID uint) error
	Encode(ctx context.Context, dto dtoreq.EncodeVideoReq) (string, error)
	CalculateDuration(ctx context.Context, dto dtoreq.CalculateVideoDurationReq) (string, error)
}

type VideoServiceImpl struct {
	minioClient     storage.Storage
	ffmpegSvc       ffmpeg.Ffmpeg
	logSvc          logger.Log
	courseSvc       courseService.CourseService
	videoRepo       repo.VideoRepo
	notificationSvc notificationService.NotificationService
	translationSvc  translations.Translator
}

func NewVideoServiceImpl(
	minioClient storage.Storage,
	ffmpegSvc ffmpeg.Ffmpeg,
	logSvc logger.Log,
	courseSvc courseService.CourseService,
	videoRepo repo.VideoRepo,
	notificationSvc notificationService.NotificationService,
	translationSvc translations.Translator,
) *VideoServiceImpl {
	return &VideoServiceImpl{
		minioClient:     minioClient,
		ffmpegSvc:       ffmpegSvc,
		logSvc:          logSvc,
		courseSvc:       courseSvc,
		videoRepo:       videoRepo,
		notificationSvc: notificationSvc,
		translationSvc:  translationSvc,
	}
}

func (svc VideoServiceImpl) AddVideo(dto *dtoreq.AddVideoToCourse) (*videoEntity.Video, error) {
	isTitleDuplicated, titleDuplicatedErr := svc.IsVideoTitleExist(dto.Title)
	if titleDuplicatedErr != nil {
		return nil, titleDuplicatedErr
	}
	if isTitleDuplicated {
		return nil, types.NewConflictError(svc.translationSvc.Translate("video.errors.title_duplicated"))
	}
	course, courseErr := svc.courseSvc.FindById(dto.CourseID)
	if courseErr != nil {
		return nil, courseErr
	}
	if course == nil {
		return nil, types.NewNotFoundError(svc.translationSvc.Translate("course.errors.not_found"))
	}
	video := &videoEntity.Video{
		Title:       dto.Title,
		IsPublished: dto.IsPublished,
		Description: dto.Description,
		AccessLevel: dto.AccessLevel,
		CourseId:    &course.ID,
		IsVerified:  false,
		Status:      videoEntity.VideoStatus_Pending,
	}
	if err := svc.videoRepo.Create(video); err != nil {
		return nil, types.NewServerError(
			"Create course throw error",
			"VideoServiceImpl.AddVideo",
			err,
		)
	}
	return video, nil
}

func (svc VideoServiceImpl) FindByTitle(title string) (*videoEntity.Video, error) {
	video, videoErr := svc.videoRepo.FetchByTitle(title)
	if videoErr != nil {
		return nil, types.NewServerError(
			"Find Video By Title Throw Error",
			"VideoServiceImpl.FetchByTitle",
			videoErr,
		)
	}
	return video, nil
}

func (svc VideoServiceImpl) IsVideoTitleExist(title string) (bool, error) {
	video, videoErr := svc.videoRepo.FetchByTitle(title)
	if videoErr != nil {
		return false, videoErr
	}
	if video == nil {
		return false, nil
	}
	return true, nil
}

func (svc VideoServiceImpl) CreateCompleteUploadVideoNotification(videoID uint) error {
	course, courseErr := svc.courseSvc.FindByVideoId(videoID)
	if courseErr != nil {
		return courseErr
	}
	if _, err := svc.notificationSvc.Create(
		notificationReqDto.NewCreateNotificationReq(
			*course.TeacherID,
			notificationEntity.NotificationType_CompleteVideoUpload,
			map[string]any{
				"videoId":     videoID,
				"courseId":    course.ID,
				"courseTitle": course.Name,
			},
		),
	); err != nil {
		return err
	}
	return nil
}

func (svc VideoServiceImpl) FindVideosByCourseID(courseID uint) ([]*videoEntity.Video, error) {
	course, courseErr := svc.courseSvc.FindById(courseID)
	if courseErr != nil {
		return nil, courseErr
	}
	if course == nil {
		return nil, types.NewNotFoundError(svc.translationSvc.Translate("course.errors.not_found"))
	}
	videos, videosErr := svc.videoRepo.FetchByCourseId(courseID)
	if videosErr != nil {
		return nil, types.NewServerError(
			"Finding videos by course id throw error",
			"VideoServiceImpl.FetchByCourseId",
			videosErr,
		)
	}
	return videos, nil
}

func (svc VideoServiceImpl) UpdateURLAndDuration(dto dtoreq.UpdateURLAndDurationVideoReq) (*videoEntity.Video, error) {
	video, videoErr := svc.FindById(dto.ID)
	if videoErr != nil {
		return nil, videoErr
	}
	if video == nil {
		return nil, types.NewNotFoundError(svc.translationSvc.Translate("video.errors.not_found"))
	}
	video.URL = dto.URL
	video.Duration = &dto.Duration
	video.Status = videoEntity.VideoStatus_Done
	if err := svc.videoRepo.Update(video); err != nil {
		return nil, types.NewServerError(
			"Error in updating the video",
			"VideoServiceImpl.UpdateVideoURL",
			err,
		)
	}
	return video, nil
}

func (svc VideoServiceImpl) FindById(id uint) (*videoEntity.Video, error) {
	video, videoErr := svc.videoRepo.FetchById(id)
	if videoErr != nil {
		return nil, types.NewServerError(
			"Error in finding video by id",
			"VideoServiceImpl.FetchById",
			videoErr,
		)
	}
	return video, nil
}

func (svc VideoServiceImpl) CalculateDuration(ctx context.Context, dto dtoreq.CalculateVideoDurationReq) (string, error) {
	log.Println("calculate duration function execute")
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
