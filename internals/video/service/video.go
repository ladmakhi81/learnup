package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoEntity "github.com/ladmakhi81/learnup/internals/video/entity"
	"github.com/ladmakhi81/learnup/internals/video/repo"
	"github.com/ladmakhi81/learnup/pkg/ffmpeg"
	"github.com/ladmakhi81/learnup/pkg/logger"
	"github.com/ladmakhi81/learnup/pkg/storage"
	"github.com/ladmakhi81/learnup/types"
	"os"
	"path"
	"path/filepath"
)

type VideoService interface {
	EncodeVideoWithObjectID(objectID string) error
	AddVideo(dto *dtoreq.AddVideoToCourse) (*videoEntity.Video, error)
	FindByTitle(title string) (*videoEntity.Video, error)
	IsVideoTitleExist(title string) (bool, error)
	FindVideosByCourseID(courseID uint) ([]*videoEntity.Video, error)
}

type VideoServiceImpl struct {
	minioClient storage.Storage
	ffmpegSvc   ffmpeg.Ffmpeg
	logSvc      logger.Log
	courseSvc   courseService.CourseService
	videoRepo   repo.VideoRepo
}

func NewVideoServiceImpl(
	minioClient storage.Storage,
	ffmpegSvc ffmpeg.Ffmpeg,
	logSvc logger.Log,
	courseSvc courseService.CourseService,
	videoRepo repo.VideoRepo,
) *VideoServiceImpl {
	return &VideoServiceImpl{
		minioClient: minioClient,
		ffmpegSvc:   ffmpegSvc,
		logSvc:      logSvc,
		courseSvc:   courseSvc,
		videoRepo:   videoRepo,
	}
}

func (svc VideoServiceImpl) EncodeVideoWithObjectID(objectID string) error {
	storeLocation, storeLocationErr := svc.encodeVideo(objectID)
	if storeLocationErr != nil {
		return storeLocationErr
	}
	if err := svc.moveVideosFromLocalToObjectStorage(storeLocation); err != nil {
		return err
	}
	if err := svc.removeUnusedVideoFiles(objectID, storeLocation); err != nil {
		return err
	}
	return nil
}

func (svc VideoServiceImpl) AddVideo(dto *dtoreq.AddVideoToCourse) (*videoEntity.Video, error) {
	isTitleDuplicated, titleDuplicatedErr := svc.IsVideoTitleExist(dto.Title)
	if titleDuplicatedErr != nil {
		return nil, titleDuplicatedErr
	}
	if isTitleDuplicated {
		return nil, types.NewConflictError("video title is duplicated")
	}
	course, courseErr := svc.courseSvc.FindById(dto.CourseID)
	if courseErr != nil {
		return nil, courseErr
	}
	if course == nil {
		return nil, types.NewNotFoundError("course not found")
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
	video, videoErr := svc.videoRepo.FindByTitle(title)
	if videoErr != nil {
		return nil, types.NewServerError(
			"Find Video By Title Throw Error",
			"VideoServiceImpl.FindByTitle",
			videoErr,
		)
	}
	return video, nil
}

func (svc VideoServiceImpl) IsVideoTitleExist(title string) (bool, error) {
	video, videoErr := svc.videoRepo.FindByTitle(title)
	if videoErr != nil {
		return false, videoErr
	}
	if video == nil {
		return false, nil
	}
	return true, nil
}

func (svc VideoServiceImpl) encodeVideo(objectID string) (string, error) {
	file, fileErr := svc.minioClient.GetFileReader(context.TODO(), "videos", objectID)
	if fileErr != nil {
		return "", types.NewServerError(
			"Error in get file from minio",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			fileErr,
		)
	}
	storeLocation, storeLocationErr := svc.ffmpegSvc.EncodeVideo(file)
	if storeLocationErr != nil {
		return "", types.NewServerError(
			"Error in encoding video file",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			storeLocationErr,
		)
	}
	return storeLocation, nil
}

func (svc VideoServiceImpl) moveVideosFromLocalToObjectStorage(storeLocation string) error {
	dirFiles, dirErr := os.ReadDir(storeLocation)
	if dirErr != nil {
		return types.NewServerError(
			"Error in finding directories of encoded files",
			"VideoServiceImpl.EncodeVideoWithObjectID",
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
			return types.NewServerError(
				"Error in finding file of directories of encoded files",
				"VideoServiceImpl.EncodeVideoWithObjectID",
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
			return types.NewServerError(
				"Error in storing encoded video into storage",
				"VideoServiceImpl.EncodeVideoWithObjectID",
				err,
			)
		}
	}
	return nil
}

func (svc VideoServiceImpl) removeUnusedVideoFiles(objectID, storeLocation string) error {
	if err := os.RemoveAll(storeLocation); err != nil {
		return types.NewServerError(
			"Error in deleting file from disk",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			err,
		)
	}
	if err := svc.minioClient.DeleteObject(context.TODO(), "videos", objectID); err != nil {
		return types.NewServerError(
			"Error in deleting file from minio",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			err,
		)
	}
	if err := svc.minioClient.DeleteObject(context.TODO(), "videos", fmt.Sprintf("%s.info", objectID)); err != nil {
		return types.NewServerError(
			"Error in deleting file from minio",
			"VideoServiceImpl.EncodeVideoWithObjectID",
			err,
		)
	}
	return nil
}

func (svc VideoServiceImpl) FindVideosByCourseID(courseID uint) ([]*videoEntity.Video, error) {
	course, courseErr := svc.courseSvc.FindById(courseID)
	if courseErr != nil {
		return nil, courseErr
	}
	if course == nil {
		return nil, types.NewNotFoundError("course not found")
	}
	videos, videosErr := svc.videoRepo.FindVideosByCourseID(courseID)
	if videosErr != nil {
		return nil, types.NewServerError(
			"Finding videos by course id throw error",
			"VideoServiceImpl.FindVideosByCourseID",
			videosErr,
		)
	}
	return videos, nil
}
