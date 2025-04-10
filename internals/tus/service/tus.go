package service

import (
	reqdto "github.com/ladmakhi81/learnup/internals/tus/dto"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/logger"
	"strconv"
)

type TusService interface {
	VideoWebhook(dto reqdto.TusWebhookDTO)
}

type TusServiceImpl struct {
	videoSvc videoService.VideoService
	logSvc   logger.Log
}

func NewTusServiceImpl(
	videoSvc videoService.VideoService,
	logSvc logger.Log,
) *TusServiceImpl {
	return &TusServiceImpl{
		videoSvc: videoSvc,
		logSvc:   logSvc,
	}
}

func (tus *TusServiceImpl) VideoWebhook(dto reqdto.TusWebhookDTO) {
	objectId, objectIdExist := dto.Event.Upload.Storage["Key"]
	courseIdParam, courseIdExist := dto.Event.Upload.MetaData["courseId"]
	courseId, courseIdErr := strconv.Atoi(courseIdParam.(string))
	if courseIdErr != nil {
		tus.logSvc.Error(logger.LogMessage{Message: "Error in converting course id"})

		return
	}
	if courseIdExist && objectIdExist {
		err := tus.videoSvc.EncodeVideoWithObjectID(uint(courseId), objectId.(string))
		if err != nil {
			tus.logSvc.Error(logger.LogMessage{
				Message: "Error happen in encoding video by ffmpeg",
				Metadata: map[string]any{
					"error": err,
					"key":   objectId,
				},
			})
		}
	}
}
