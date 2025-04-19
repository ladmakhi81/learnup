package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	questionService "github.com/ladmakhi81/learnup/internals/question/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
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
	GetQuestions(options GetQuestionOptions) ([]*entities.Question, error)
	GetQuestionCount(courseId *uint) (int, error)
}

type TeacherQuestionServiceImpl struct {
	questionSvc    questionService.QuestionService
	translationSvc contracts.Translator
	userSvc        userService.UserSvc
	courseSvc      courseService.CourseService
}

func NewTeacherQuestionServiceImpl(
	questionSvc questionService.QuestionService,
	translationSvc contracts.Translator,
	userSvc userService.UserSvc,
	courseSvc courseService.CourseService,
) *TeacherQuestionServiceImpl {
	return &TeacherQuestionServiceImpl{
		questionSvc:    questionSvc,
		translationSvc: translationSvc,
		userSvc:        userSvc,
		courseSvc:      courseSvc,
	}
}

func (svc TeacherQuestionServiceImpl) GetQuestions(options GetQuestionOptions) ([]*entities.Question, error) {
	teacher, teacherErr := svc.userSvc.FindById(options.TeacherID)
	if teacherErr != nil {
		return nil, teacherErr
	}
	if teacher == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	if options.CourseID != nil {
		course, courseErr := svc.courseSvc.FindById(*options.CourseID)
		if courseErr != nil {
			return nil, courseErr
		}
		if course == nil {
			return nil, types.NewNotFoundError(
				svc.translationSvc.Translate("course.errors.not_found"),
			)
		}
		if *course.TeacherID != teacher.ID {
			return nil, types.NewForbiddenAccessError(
				svc.translationSvc.Translate("common.errors.forbidden_access"),
			)
		}
	}
	questions, questionsErr := svc.questionSvc.GetPageable(
		options.CourseID,
		options.Page,
		options.PageSize,
	)
	if questionsErr != nil {
		return nil, questionsErr
	}
	return questions, nil
}

func (svc TeacherQuestionServiceImpl) GetQuestionCount(courseId *uint) (int, error) {
	count, countErr := svc.questionSvc.GetCount(courseId)
	if countErr != nil {
		return 0, countErr
	}
	return count, nil
}
