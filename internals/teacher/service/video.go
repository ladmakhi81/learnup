package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	notificationEntity "github.com/ladmakhi81/learnup/db/entities"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	dtoreq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	videoRepository "github.com/ladmakhi81/learnup/internals/video/repo"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherVideoService interface {
	AddVideo(dto dtoreq.AddVideoToCourseReq) (*entities.Video, error)
}

type TeacherVideoServiceImpl struct {
	videoSvc       videoService.VideoService
	translationSvc contracts.Translator
	courseSvc      courseService.CourseService
	videoRepo      videoRepository.VideoRepo
}

func NewTeacherVideoServiceImpl(
	videoSvc videoService.VideoService,
	translationSvc contracts.Translator,
	courseSvc courseService.CourseService,
	videoRepo videoRepository.VideoRepo,
) *TeacherVideoServiceImpl {
	return &TeacherVideoServiceImpl{
		videoSvc:       videoSvc,
		translationSvc: translationSvc,
		courseSvc:      courseSvc,
		videoRepo:      videoRepo,
	}
}

func (svc TeacherVideoServiceImpl) AddVideo(dto dtoreq.AddVideoToCourseReq) (*entities.Video, error) {
	isTitleDuplicated, titleDuplicatedErr := svc.videoSvc.IsVideoTitleExist(dto.Title)
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
	video := &notificationEntity.Video{
		Title:       dto.Title,
		IsPublished: dto.IsPublished,
		Description: dto.Description,
		AccessLevel: dto.AccessLevel,
		CourseId:    &course.ID,
		IsVerified:  false,
		Status:      notificationEntity.VideoStatus_Pending,
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
