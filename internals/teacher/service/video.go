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
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewTeacherVideoServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *TeacherVideoServiceImpl {
	return &TeacherVideoServiceImpl{
		translationSvc: translationSvc,
		unitOfWork:     unitOfWork,
	}
}

func (svc TeacherVideoServiceImpl) AddVideo(dto dtoreq.AddVideoToCourseReq) (*entities2.Video, error) {
	const operationName = "TeacherVideoServiceImpl.AddVideo"
	isTitleDuplicated, err := svc.unitOfWork.VideoRepo().Exist(map[string]any{
		"title": dto.Title,
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in checking existence of video title",
			operationName,
			err,
		)
	}
	if isTitleDuplicated {
		return nil, types.NewConflictError(svc.translationSvc.Translate("video.errors.title_duplicated"))
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching course by id",
			operationName,
			err,
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
	if err := svc.unitOfWork.VideoRepo().Create(video); err != nil {
		return nil, types.NewServerError(
			"Create course throw error",
			operationName,
			err,
		)
	}
	return video, nil
}
