package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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
}

type VideoServiceImpl struct {
	minioClient storage.Storage
	ffmpegSvc   ffmpeg.Ffmpeg
	logSvc      logger.Log
}

func NewVideoServiceImpl(
	minioClient storage.Storage,
	ffmpegSvc ffmpeg.Ffmpeg,
	logSvc logger.Log,
) *VideoServiceImpl {
	return &VideoServiceImpl{
		minioClient: minioClient,
		ffmpegSvc:   ffmpegSvc,
		logSvc:      logSvc,
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
