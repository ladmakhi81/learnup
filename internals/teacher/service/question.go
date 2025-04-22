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
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewTeacherQuestionServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *TeacherQuestionServiceImpl {
	return &TeacherQuestionServiceImpl{
		translationSvc: translationSvc,
		unitOfWork:     unitOfWork,
	}
}

func (svc TeacherQuestionServiceImpl) GetQuestions(options GetQuestionOptions) ([]*entities.Question, int, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, 0, txErr
	}
	teacher, teacherErr := tx.UserRepo().GetByID(options.TeacherID, nil)
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
		course, courseErr := tx.CourseRepo().GetByID(*options.CourseID, nil)
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
	questionCondition := make(map[string]any)
	if options.CourseID != nil {
		questionCondition["course_id"] = *options.CourseID
	}
	questions, count, questionsErr := tx.QuestionRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset:     &options.Page,
			Limit:      &options.PageSize,
			Conditions: questionCondition,
			Relations:  []string{"User", "Course", "Video"},
		},
	)
	if questionsErr != nil {
		return nil, 0, questionsErr
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return questions, count, nil
}
