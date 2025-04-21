package service

import (
	"context"
	"fmt"
	reqdto "github.com/ladmakhi81/learnup/internals/tus/dto"
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/internals/video/workflow"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/pkg/temporal"
	"github.com/ladmakhi81/learnup/utils"
	"strconv"
)

type TusService interface {
	VideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO)
	AddCourseVideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO)
	AddIntroductionVideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO)
}

type TusServiceImpl struct {
	videoSvc         videoService.VideoService
	logSvc           contracts.Log
	temporalSvc      contracts.Temporal
	videoWorkflowSvc workflow.VideoWorkflow
}

func NewTusServiceImpl(
	videoSvc videoService.VideoService,
	logSvc contracts.Log,
	temporalSvc contracts.Temporal,
	videoWorkflowSvc workflow.VideoWorkflow,
) *TusServiceImpl {
	return &TusServiceImpl{
		videoSvc:         videoSvc,
		logSvc:           logSvc,
		temporalSvc:      temporalSvc,
		videoWorkflowSvc: videoWorkflowSvc,
	}
}

func (tus TusServiceImpl) VideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO) {
	fmt.Println(1, fmt.Sprintf("%v", reqdto.TusActionType_AddIntroductionVideo) == dto.Event.Upload.MetaData["actionType"], fmt.Sprintf("--%s--", dto.Event.Upload.MetaData["actionType"]))
	switch dto.Event.Upload.MetaData["actionType"] {
	case utils.ToString(reqdto.TusActionType_NewCourseVideo):
		tus.AddCourseVideoWebhook(ctx, dto)
	case utils.ToString(reqdto.TusActionType_AddIntroductionVideo):
		tus.AddIntroductionVideoWebhook(ctx, dto)
	}
}

func (tus TusServiceImpl) AddCourseVideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO) {
	objectId, objectIdExist := dto.Event.Upload.Storage["Key"]
	courseIdParam, courseIdExist := dto.Event.Upload.MetaData["courseId"]
	videoIdParam, videoIdExist := dto.Event.Upload.MetaData["videoId"]
	if objectIdExist && courseIdExist && videoIdExist {
		courseId, courseIdErr := strconv.Atoi(courseIdParam.(string))
		if courseIdErr != nil {
			tus.logSvc.Error(dtos.LogMessage{Message: "Error in converting course id"})
			return
		}
		videoId, videoIdErr := strconv.Atoi(videoIdParam.(string))
		if videoIdErr != nil {
			tus.logSvc.Error(dtos.LogMessage{Message: "Error in converting video id"})
			return
		}
		workflowDto := dtoreq.AddNewCourseVideoWorkflowReq{
			CourseID: uint(courseId),
			ObjectID: objectId.(string),
			VideoID:  uint(videoId),
		}
		workflowErr := tus.temporalSvc.ExecuteWorker(
			ctx,
			temporal.ADD_NEW_COURSE_VIDEO_QUEUE,
			tus.videoWorkflowSvc.AddNewCourseVideoWorkflow,
			workflowDto,
		)

		if workflowErr != nil {
			tus.logSvc.Error(dtos.LogMessage{
				Message: "Error happen in workflow of video workflow service",
				Metadata: map[string]any{
					"error": workflowErr,
					"key":   objectId,
				},
			})
		}
	}

	return
}

func (tus TusServiceImpl) AddIntroductionVideoWebhook(ctx context.Context, dto reqdto.TusWebhookDTO) {
	objectId, objectIdExist := dto.Event.Upload.Storage["Key"]
	courseIdParam, courseIdExist := dto.Event.Upload.MetaData["courseId"]
	if objectIdExist && courseIdExist {
		courseId, courseIdErr := strconv.Atoi(courseIdParam.(string))
		if courseIdErr != nil {
			tus.logSvc.Error(dtos.LogMessage{Message: "Error in converting course id"})
			return
		}
		workflowDto := dtoreq.AddIntroductionVideoWorkflowReq{
			CourseId: uint(courseId),
			ObjectId: objectId.(string),
		}
		workflowErr := tus.temporalSvc.ExecuteWorker(
			ctx,
			temporal.SET_INTRODUCTION_COURSE_QUEUE,
			tus.videoWorkflowSvc.AddIntroductionVideoWorkflow,
			workflowDto,
		)
		if workflowErr != nil {
			tus.logSvc.Error(dtos.LogMessage{
				Message: "Error happen in workflow of video workflow service",
				Metadata: map[string]any{
					"error": workflowErr,
					"key":   objectId,
				},
			})
		}
	}
}
