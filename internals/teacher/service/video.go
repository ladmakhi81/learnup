package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherVideoService interface {
	AddVideo(dto dtoreq.AddVideoToCourseReq) (*entities2.Video, error)
}

type TeacherVideoServiceImpl struct {
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewTeacherVideoServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *TeacherVideoServiceImpl {
	return &TeacherVideoServiceImpl{
		translationSvc: translationSvc,
		repo:           repo,
	}
}

func (svc TeacherVideoServiceImpl) AddVideo(dto dtoreq.AddVideoToCourseReq) (*entities2.Video, error) {
	isTitleDuplicated, titleDuplicatedErr := svc.repo.VideoRepo.Exist(map[string]any{
		"title": dto.Title,
	})
	if titleDuplicatedErr != nil {
		return nil, types.NewServerError(
			"Error in checking existence of video title",
			"TeacherVideoServiceImpl.AddVideo",
			titleDuplicatedErr,
		)
	}
	if isTitleDuplicated {
		return nil, types.NewConflictError(svc.translationSvc.Translate("video.errors.title_duplicated"))
	}
	course, courseErr := svc.repo.CourseRepo.GetByID(dto.CourseID)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in fetching course by id",
			"TeacherVideoServiceImpl.AddVideo",
			courseErr,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(svc.translationSvc.Translate("course.errors.not_found"))
	}
	video := &entities2.Video{
		Title:       dto.Title,
		IsPublished: dto.IsPublished,
		Description: dto.Description,
		AccessLevel: dto.AccessLevel,
		CourseId:    &course.ID,
		IsVerified:  false,
		Status:      entities2.VideoStatus_Pending,
	}
	if err := svc.repo.VideoRepo.Create(video); err != nil {
		return nil, types.NewServerError(
			"Create course throw error",
			"VideoServiceImpl.AddVideo",
			err,
		)
	}
	return video, nil
}
