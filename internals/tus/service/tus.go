package service

import (
	reqdto "github.com/ladmakhi81/learnup/internals/tus/dto"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/logger"
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
	switch dto.Type {
	case reqdto.TusHookType_PostFinish:
		if objectId, objectIdExist := dto.Event.Upload.Storage["Key"]; objectIdExist {
			go func() {
				err := tus.videoSvc.EncodeVideoWithObjectID(objectId.(string))
				if err != nil {
					tus.logSvc.Error(logger.LogMessage{
						Message: "Error happen in encoding video by ffmpeg",
						Metadata: map[string]any{
							"error": err,
							"key":   objectId,
						},
					})
				}
			}()
		}
	}
}
