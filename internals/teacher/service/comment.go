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
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, 0, txErr
	}
	teacherClaim := authContext.(*types.TokenClaim)
	teacher, teacherErr := tx.UserRepo().GetByID(teacherClaim.UserID, nil)
	if teacherErr != nil {
		return nil, 0, types.NewServerError(
			"Error in finding teacher by id",
			"TeacherCommentServiceImpl.GetPageableCommentByCourseId",
			teacherErr,
		)
	}
	if teacher == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	course, courseErr := tx.CourseRepo().GetByID(courseId, nil)
	if courseErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching course detail",
			"TeacherCommentServiceImpl.GetPageableCommentByCourseId",
			courseErr,
		)
	}
	if course == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	comments, count, commentsErr := tx.CommentRepo().GetPaginated(repositories.GetPaginatedOptions{
		Limit:     &pageSize,
		Offset:    &page,
		Relations: []string{"User", "Course"},
		Conditions: map[string]any{
			"course_id": courseId,
			"user_id":   teacher.ID,
		},
	})
	if commentsErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching comments",
			"TeacherCommentServiceImpl.GetPageableCommentByCourseId",
			commentsErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, txErr
	}
	return comments, count, nil
}
