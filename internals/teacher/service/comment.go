package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/db/repositories"
	"github.com/ladmakhi81/learnup/shared/types"
)

type TeacherCommentService interface {
	GetPageableCommentByCourseId(teacher *entities.User, courseId uint, page, pageSize int) ([]*entities.Comment, int, error)
}

type teacherCommentService struct {
	unitOfWork db.UnitOfWork
}

func NewTeacherCommentSvc(unitOfWork db.UnitOfWork) TeacherCommentService {
	return &teacherCommentService{unitOfWork: unitOfWork}
}

func (svc teacherCommentService) GetPageableCommentByCourseId(teacher *entities.User, courseId uint, page, pageSize int) ([]*entities.Comment, int, error) {
	const operationName = "teacherCommentService.GetPageableCommentByCourseId"
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
