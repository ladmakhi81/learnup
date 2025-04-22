package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	questionDtoReq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type QuestionAnswerService interface {
	Create(dto questionDtoReq.AnswerQuestionReq) (*entities.QuestionAnswer, error)
	GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error)
}

type QuestionAnswerServiceImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewQuestionAnswerServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *QuestionAnswerServiceImpl {
	return &QuestionAnswerServiceImpl{
		translationSvc: translationSvc,
		unitOfWork:     unitOfWork,
	}
}

func (svc QuestionAnswerServiceImpl) Create(dto questionDtoReq.AnswerQuestionReq) (*entities.QuestionAnswer, error) {
	const operationName = "QuestionAnswerServiceImpl.Create"
	sender, err := svc.unitOfWork.UserRepo().GetByID(dto.SenderID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching user",
			operationName,
			err,
		)
	}
	if sender == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	question, err := svc.unitOfWork.QuestionRepo().GetByID(dto.QuestionID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching question by id",
			operationName,
			err,
		)
	}
	if question == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("question.errors.not_found"),
		)
	}
	if question.IsClosed {
		return nil, types.NewBadRequestError(
			svc.translationSvc.Translate("question.errors.closed"),
		)
	}
	answer := &entities.QuestionAnswer{
		QuestionID: dto.QuestionID,
		Content:    dto.Content,
		SenderID:   sender.ID,
	}
	if err := svc.unitOfWork.AnswerRepo().Create(answer); err != nil {
		return nil, types.NewServerError(
			"Error in creating answer",
			operationName,
			err,
		)
	}
	return answer, nil
}

func (svc QuestionAnswerServiceImpl) GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error) {
	const operationName = "QuestionAnswerServiceImpl.GetQuestionAnswers"
	question, err := svc.unitOfWork.QuestionRepo().GetByID(questionID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching question by id",
			operationName,
			err,
		)
	}
	if question == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("question.errors.not_found"),
		)
	}
	answers, err := svc.unitOfWork.AnswerRepo().GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"question_id": questionID,
		},
		Relations: []string{"Sender"},
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching answers",
			operationName,
			err,
		)
	}
	return answers, nil
}
