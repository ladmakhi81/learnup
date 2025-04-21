package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type GetQuestionOptions struct {
	TeacherID uint
	CourseID  *uint
	Page      int
	PageSize  int
}

type TeacherQuestionService interface {
	GetQuestions(options GetQuestionOptions) ([]*entities.Question, int, error)
}

type TeacherQuestionServiceImpl struct {
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewTeacherQuestionServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *TeacherQuestionServiceImpl {
	return &TeacherQuestionServiceImpl{
		translationSvc: translationSvc,
		repo:           repo,
	}
}

func (svc TeacherQuestionServiceImpl) GetQuestions(options GetQuestionOptions) ([]*entities.Question, int, error) {
	teacher, teacherErr := svc.repo.UserRepo.GetByID(options.TeacherID)
	if teacherErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching teacher by id",
			"TeacherQuestionServiceImpl.GetQuestions",
			teacherErr,
		)
	}
	if teacher == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	if options.CourseID != nil {
		course, courseErr := svc.repo.CourseRepo.GetByID(*options.CourseID)
		if courseErr != nil {
			return nil, 0, types.NewServerError(
				"Error in fetching course by id",
				"TeacherQuestionServiceImpl.GetQuestions",
				courseErr,
			)
		}
		if course == nil {
			return nil, 0, types.NewNotFoundError(
				svc.translationSvc.Translate("course.errors.not_found"),
			)
		}
		if *course.TeacherID != teacher.ID {
			return nil, 0, types.NewForbiddenAccessError(
				svc.translationSvc.Translate("common.errors.forbidden_access"),
			)
		}
	}
	questions, count, questionsErr := svc.repo.QuestionRepo.GetPaginated(
		repositories.GetPaginatedOptions{
			Offset: &options.Page,
			Limit:  &options.PageSize,
			Conditions: map[string]any{
				"course_id": options.CourseID,
			},
		},
	)
	if questionsErr != nil {
		return nil, 0, questionsErr
	}
	return questions, count, nil
}
