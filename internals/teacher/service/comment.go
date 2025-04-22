package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherCommentService interface {
	GetPageableCommentByCourseId(authContext any, courseId uint, page, pageSize int) ([]*entities.Comment, int, error)
}

type TeacherCommentServiceImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewTeacherCommentServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *TeacherCommentServiceImpl {
	return &TeacherCommentServiceImpl{
		translationSvc: translationSvc,
		unitOfWork:     unitOfWork,
	}
}

func (svc TeacherCommentServiceImpl) GetPageableCommentByCourseId(authContext any, courseId uint, page, pageSize int) ([]*entities.Comment, int, error) {
	const operationName = "TeacherCommentServiceImpl.GetPageableCommentByCourseId"
	teacherClaim := authContext.(*types.TokenClaim)
	teacher, err := svc.unitOfWork.UserRepo().GetByID(teacherClaim.UserID, nil)
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in finding teacher by id",
			operationName,
			err,
		)
	}
	if teacher == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(courseId, nil)
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching course detail",
			operationName,
			err,
		)
	}
	if course == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
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
		return nil, 0, types.NewServerError(
			"Error in fetching comments",
			operationName,
			err,
		)
	}
	return comments, count, nil
}
