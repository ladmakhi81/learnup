package workflow

import (
	videoDtoReq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoEntity "github.com/ladmakhi81/learnup/internals/video/entity"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/temporal"
	"go.temporal.io/sdk/workflow"
)

type VideoWorkflow interface {
	VideoCourseWorkflow(ctx workflow.Context, dto videoDtoReq.VideoCourseWorkflowDto) error
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

func (svc VideoWorkflowImpl) VideoCourseWorkflow(ctx workflow.Context, dto videoDtoReq.VideoCourseWorkflowDto) error {
	// calculate duration
	calculateDurationDto := videoDtoReq.CalculateVideoDurationReq{
		ObjectId: dto.ObjectID,
	}
	var videoDuration string
	calculateDurationErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.CalculateDuration, calculateDurationDto, &videoDuration)
	if calculateDurationErr != nil {
		return calculateDurationErr
	}
	// encode
	var videoURL string
	encodeVideoDto := videoDtoReq.EncodeVideoReq{
		ObjectId: dto.ObjectID,
	}
	encodeErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.Encode, encodeVideoDto, &videoURL)
	if encodeErr != nil {
		return encodeErr
	}
	// update url and duration
	var video *videoEntity.Video
	updateVideoDto := videoDtoReq.UpdateURLAndDurationVideoReq{
		Duration: videoDuration,
		URL:      videoURL,
		ID:       dto.VideoID,
	}
	updateErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.UpdateURLAndDuration, updateVideoDto, &video)
	if updateErr != nil {
		return updateErr
	}

	// teacher notification
	teacherNotificationErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.CreateCompleteUploadVideoNotification, video.ID, nil)
	if teacherNotificationErr != nil {
		return teacherNotificationErr
	}
	return nil
}
