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
	const operationName = "TeacherQuestionServiceImpl.GetQuestions"
	teacher, err := svc.unitOfWork.UserRepo().GetByID(options.TeacherID, nil)
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching teacher by id",
			operationName,
			err,
		)
	}
	if teacher == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	if options.CourseID != nil {
		course, err := svc.unitOfWork.CourseRepo().GetByID(*options.CourseID, nil)
		if err != nil {
			return nil, 0, types.NewServerError(
				"Error in fetching course by id",
				operationName,
				err,
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
	questions, count, err := svc.unitOfWork.QuestionRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset:     &options.Page,
			Limit:      &options.PageSize,
			Conditions: questionCondition,
			Relations:  []string{"User", "Course", "Video"},
		},
	)
	if err != nil {
		return nil, 0, err
	}
	return questions, count, nil
}
