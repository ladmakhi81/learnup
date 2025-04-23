package workflow

import (
	courseDtoReq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	videoDtoReq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	videoEntity "github.com/ladmakhi81/learnup/shared/db/entities"
	"go.temporal.io/sdk/workflow"
)

type VideoWorkflow interface {
	AddNewCourseVideoWorkflow(ctx workflow.Context, dto videoDtoReq.AddNewCourseVideoWorkflowReqDto) error
	AddIntroductionVideoWorkflow(ctx workflow.Context, dto videoDtoReq.AddIntroductionVideoWorkflowReqDto) error
}

type VideoWorkflowImpl struct {
	videoSvc    videoService.VideoService
	temporalSvc contracts.Temporal
	courseSvc   courseService.CourseService
}

func NewVideoWorkflowImpl(
	videoSvc videoService.VideoService,
	temporal contracts.Temporal,
	courseSvc courseService.CourseService,
) *VideoWorkflowImpl {
	return &VideoWorkflowImpl{
		videoSvc:    videoSvc,
		temporalSvc: temporal,
		courseSvc:   courseSvc,
	}
}

func (svc VideoWorkflowImpl) AddNewCourseVideoWorkflow(ctx workflow.Context, dto videoDtoReq.AddNewCourseVideoWorkflowReqDto) error {
	// calculate duration
	calculateDurationDto := videoDtoReq.CalculateVideoDurationReqDto{
		ObjectId: dto.ObjectID,
	}
	var videoDuration string
	calculateDurationErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.CalculateDuration, calculateDurationDto, &videoDuration)
	if calculateDurationErr != nil {
		return calculateDurationErr
	}
	// encode
	var videoURL string
	encodeVideoDto := videoDtoReq.EncodeVideoReqDto{
		ObjectId: dto.ObjectID,
	}
	encodeErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.Encode, encodeVideoDto, &videoURL)
	if encodeErr != nil {
		return encodeErr
	}
	// update url and duration
	var video *videoEntity.Video
	updateVideoDto := videoDtoReq.UpdateURLAndDurationVideoReqDto{
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

func (svc VideoWorkflowImpl) AddIntroductionVideoWorkflow(ctx workflow.Context, dto videoDtoReq.AddIntroductionVideoWorkflowReqDto) error {
	// encode
	var videoURL string
	encodeVideoDto := videoDtoReq.EncodeVideoReqDto{
		ObjectId: dto.ObjectId,
	}
	encodeErr := svc.temporalSvc.ExecuteTask(ctx, svc.videoSvc.Encode, encodeVideoDto, &videoURL)
	if encodeErr != nil {
		return encodeErr
	}
	// update video introduction url
	updateIntroductionUrlDto := courseDtoReq.UpdateIntroductionURLReqDto{
		URL:      videoURL,
		CourseId: dto.CourseId,
	}
	updateErr := svc.temporalSvc.ExecuteTask(ctx, svc.courseSvc.UpdateIntroductionURL, updateIntroductionUrlDto, nil)
	if updateErr != nil {
		return updateErr
	}
	// create notification
	notificationErr := svc.temporalSvc.ExecuteTask(ctx, svc.courseSvc.CreateCompleteIntroductionVideoNotification, dto.CourseId, nil)
	if notificationErr != nil {
		return notificationErr
	}
	return nil
}
