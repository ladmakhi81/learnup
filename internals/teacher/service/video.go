package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/internals/db"
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	videoError "github.com/ladmakhi81/learnup/internals/video/error"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherVideoService interface {
	AddVideo(dto dtoreq.AddVideoToCourseReqDto) (*entities2.Video, error)
}

type teacherVideoService struct {
	unitOfWork db.UnitOfWork
}

func NewTeacherVideoSvc(unitOfWork db.UnitOfWork) TeacherVideoService {
	return &teacherVideoService{unitOfWork: unitOfWork}
}

func (svc teacherVideoService) AddVideo(dto dtoreq.AddVideoToCourseReqDto) (*entities2.Video, error) {
	const operationName = "teacherVideoService.AddVideo"
	isTitleDuplicated, err := svc.unitOfWork.VideoRepo().Exist(map[string]any{"title": dto.Title})
	if err != nil {
		return nil, types.NewServerError("Error in checking existence of video title", operationName, err)
	}
	if isTitleDuplicated {
		return nil, videoError.Video_TitleDuplicated
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching course by id", operationName, err)
	}
	if course == nil {
		return nil, courseError.Course_NotFound
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
		return nil, types.NewServerError("Create course throw error", operationName, err)
	}
	return video, nil
}
