package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	dtoreq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type QuestionService interface {
	Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error)
	GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, int, error)
}

type QuestionServiceImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewQuestionServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *QuestionServiceImpl {
	return &QuestionServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc QuestionServiceImpl) Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error) {
	const operationName = "QuestionServiceImpl.Create"
	sender, err := svc.unitOfWork.UserRepo().GetByID(dto.UserID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching sender data",
			operationName,
			err,
		)
	}
	if sender == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching course by id",
			operationName,
			err,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	question := &entities.Question{
		UserID:   sender.ID,
		CourseID: course.ID,
		Content:  dto.Content,
		Priority: dto.Priority,
	}
	if dto.VideoID != nil {
		video, err := svc.unitOfWork.VideoRepo().GetByID(*dto.VideoID, nil)
		if err != nil {
			return nil, types.NewServerError(
				"Error in fetching video",
				operationName,
				err,
			)
		}
		if video == nil {
			return nil, types.NewNotFoundError(
				svc.translationSvc.Translate("video.errors.not_found"),
			)
		}
		question.VideoID = &video.ID
	}
	if err := svc.unitOfWork.QuestionRepo().Create(question); err != nil {
		return nil, types.NewServerError(
			"Error in creating question",
			operationName,
			err,
		)
	}
	// TODO: notification system
	// send notification for teacher that we have new question
	return question, nil
}

func (svc QuestionServiceImpl) GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, int, error) {
	const operationName = "QuestionServiceImpl.GetPageable"
	questions, count, err := svc.unitOfWork.QuestionRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Limit:  &pageSize,
			Offset: &page,
			Conditions: map[string]any{
				"course_id": courseId,
			},
			Relations: []string{
				"User",
				"Course",
				"Video",
			},
		})
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching related questions",
			operationName,
			err,
		)
	}
	return questions, count, nil
}
