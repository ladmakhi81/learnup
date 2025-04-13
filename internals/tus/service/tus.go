package service

import (
	"context"
	reqdto "github.com/ladmakhi81/learnup/internals/tus/dto"
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/internals/video/workflow"
	"github.com/ladmakhi81/learnup/pkg/logger"
	"github.com/ladmakhi81/learnup/pkg/temporal"
	"strconv"
)

type TusService interface {
	VideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO)
}

type TusServiceImpl struct {
	videoSvc         videoService.VideoService
	logSvc           logger.Log
	temporalSvc      temporal.Temporal
	videoWorkflowSvc workflow.VideoWorkflow
}

func NewTusServiceImpl(
	videoSvc videoService.VideoService,
	logSvc logger.Log,
	temporalSvc temporal.Temporal,
	videoWorkflowSvc workflow.VideoWorkflow,
) *TusServiceImpl {
	return &TusServiceImpl{
		videoSvc:         videoSvc,
		logSvc:           logSvc,
		temporalSvc:      temporalSvc,
		videoWorkflowSvc: videoWorkflowSvc,
	}
}

func (tus *TusServiceImpl) VideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO) {
	objectId, objectIdExist := dto.Event.Upload.Storage["Key"]
	courseIdParam, courseIdExist := dto.Event.Upload.MetaData["courseId"]
	courseId, courseIdErr := strconv.Atoi(courseIdParam.(string))
	if courseIdErr != nil {
		tus.logSvc.Error(logger.LogMessage{Message: "Error in converting course id"})
		return
	}
	if objectIdExist && courseIdExist {
		workflowDto := dtoreq.VideoCourseWorkflowDto{
			CourseID: uint(courseId),
			ObjectID: objectId.(string),
		}
		workflowErr := tus.temporalSvc.ExecuteWorker(
			ctx,
			temporal.COURSE_VIDEO_QUEUE,
			tus.videoWorkflowSvc.VideoCourseWorkflow,
			workflowDto,
		)

		if workflowErr != nil {
			tus.logSvc.Error(logger.LogMessage{
				Message: "Error happen in workflow of video workflow svc",
				Metadata: map[string]any{
					"error": workflowErr,
					"key":   objectId,
				},
			})
		}
	}

	return
}
