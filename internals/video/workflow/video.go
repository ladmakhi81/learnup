package workflow

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/temporal"
	"go.temporal.io/sdk/workflow"
	"log"
)

type VideoWorkflow interface {
	VideoCourseWorkflow(ctx workflow.Context, dto dtoreq.VideoCourseWorkflowDto) error
}

type VideoWorkflowImpl struct {
	videoSvc    videoService.VideoService
	temporalSvc temporal.Temporal
}

func NewVideoWorkflowImpl(
	videoSvc videoService.VideoService,
	temporal temporal.Temporal,
) *VideoWorkflowImpl {
	return &VideoWorkflowImpl{
		videoSvc:    videoSvc,
		temporalSvc: temporal,
	}
}

func (svc VideoWorkflowImpl) VideoCourseWorkflow(ctx workflow.Context, dto dtoreq.VideoCourseWorkflowDto) error {
	calculateDurationDto := dtoreq.CalculateVideoDurationReq{
		ObjectId: dto.ObjectID,
	}
	var videoDuration string
	calculateDurationErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.CalculateDuration, calculateDurationDto, &videoDuration)
	if calculateDurationErr != nil {
		return calculateDurationErr
	}
	attachSubtitleErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.AttachSubtitle, map[string]any{}, nil)
	if attachSubtitleErr != nil {
		return attachSubtitleErr
	}
	encodeErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.Encode, map[string]any{}, nil)
	if encodeErr != nil {
		return encodeErr
	}

	log.Println("video duration: ", videoDuration)
	return nil
}
