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
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	isTitleDuplicated, titleDuplicatedErr := tx.VideoRepo().Exist(map[string]any{
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
	course, courseErr := tx.CourseRepo().GetByID(dto.CourseID, nil)
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
	if err := tx.VideoRepo().Create(video); err != nil {
		return nil, types.NewServerError(
			"Create course throw error",
			"VideoServiceImpl.AddVideo",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return video, nil
}
