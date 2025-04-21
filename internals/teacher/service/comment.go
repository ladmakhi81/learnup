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
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewTeacherCommentServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *TeacherCommentServiceImpl {
	return &TeacherCommentServiceImpl{
		translationSvc: translationSvc,
		repo:           repo,
	}
}

func (svc TeacherCommentServiceImpl) GetPageableCommentByCourseId(authContext any, courseId uint, page, pageSize int) ([]*entities.Comment, int, error) {
	teacherClaim := authContext.(*types.TokenClaim)
	teacher, teacherErr := svc.repo.UserRepo.GetByID(teacherClaim.UserID)
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
	course, courseErr := svc.repo.CourseRepo.GetByID(courseId)
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
	comments, count, commentsErr := svc.repo.CommentRepo.GetPaginated(repositories.GetPaginatedOptions{
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
	return comments, count, nil
}
