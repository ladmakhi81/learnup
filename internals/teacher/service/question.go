package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/types"
)

type GetQuestionOptions struct {
	CourseID *uint
	Page     int
	PageSize int
}

type TeacherQuestionService interface {
	GetQuestions(teacher *entities.User, options GetQuestionOptions) ([]*entities.Question, int, error)
}

type teacherQuestionService struct {
	unitOfWork db.UnitOfWork
}

func NewTeacherQuestionSvc(unitOfWork db.UnitOfWork) TeacherQuestionService {
	return &teacherQuestionService{unitOfWork: unitOfWork}
}

func (svc teacherQuestionService) GetQuestions(teacher *entities.User, options GetQuestionOptions) ([]*entities.Question, int, error) {
	const operationName = "teacherQuestionService.GetQuestions"
	if options.CourseID != nil {
		course, err := svc.unitOfWork.CourseRepo().GetByID(*options.CourseID, nil)
		if err != nil {
			return nil, 0, types.NewServerError("Error in fetching course by id", operationName, err)
		}
		if course == nil {
			return nil, 0, courseError.Course_NotFound
		}
		if !course.IsTeacher(teacher.ID) {
			return nil, 0, courseError.Course_ForbiddenAccess
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
