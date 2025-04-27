package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	forumError "github.com/ladmakhi81/learnup/internals/forum/error"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/types"
)

type ForumService interface {
	GetForumByCourseID(courseID uint) (*entities.CourseForum, error)
}

type forumService struct {
	unitOfWork db.UnitOfWork
}

func NewForumService(unitOfWork db.UnitOfWork) ForumService {
	return &forumService{
		unitOfWork: unitOfWork,
	}
}

func (svc forumService) GetForumByCourseID(courseID uint) (*entities.CourseForum, error) {
	const operationName = "forumService.GetForumByCourseID"
	course, err := svc.unitOfWork.CourseRepo().GetByID(courseID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in getting course forum", operationName, err)
	}
	if course == nil {
		return nil, courseError.Course_NotFound
	}
	forum, err := svc.unitOfWork.CourseForumRepo().GetOne(
		map[string]any{"course_id": courseID},
		[]string{
			"Course",
			"Teacher",
			"Course.Participants",
			"Course.Participants.Student",
		},
	)
	if err != nil {
		return nil, types.NewServerError("Error in getting forum by course id", operationName, err)
	}
	if forum == nil {
		return nil, forumError.Forum_NotFound
	}
	return forum, nil
}
