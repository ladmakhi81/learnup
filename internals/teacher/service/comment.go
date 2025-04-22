package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	userError "github.com/ladmakhi81/learnup/internals/user/error"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherCommentService interface {
	GetPageableCommentByCourseId(authContext any, courseId uint, page, pageSize int) ([]*entities.Comment, int, error)
}

type teacherCommentService struct {
	unitOfWork db.UnitOfWork
}

func NewTeacherCommentSvc(unitOfWork db.UnitOfWork) TeacherCommentService {
	return &teacherCommentService{unitOfWork: unitOfWork}
}

func (svc teacherCommentService) GetPageableCommentByCourseId(authContext any, courseId uint, page, pageSize int) ([]*entities.Comment, int, error) {
	const operationName = "teacherCommentService.GetPageableCommentByCourseId"
	teacherClaim := authContext.(*types.TokenClaim)
	teacher, err := svc.unitOfWork.UserRepo().GetByID(teacherClaim.UserID, nil)
	if err != nil {
		return nil, 0, types.NewServerError("Error in finding teacher by id", operationName, err)
	}
	if teacher == nil {
		return nil, 0, userError.User_TeacherNotFound
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(courseId, nil)
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching course detail", operationName, err)
	}
	if course == nil {
		return nil, 0, courseError.Course_NotFound
	}
	comments, count, err := svc.unitOfWork.CommentRepo().GetPaginated(repositories.GetPaginatedOptions{
		Limit:     &pageSize,
		Offset:    &page,
		Relations: []string{"User", "Course"},
		Conditions: map[string]any{
			"course_id": courseId,
			"user_id":   teacher.ID,
		},
	})
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching comments", operationName, err)
	}
	return comments, count, nil
}
