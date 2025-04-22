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
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	sender, senderErr := tx.UserRepo().GetByID(dto.SenderID, nil)
	if senderErr != nil {
		return nil, types.NewServerError(
			"Error in fetching user",
			"QuestionAnswerServiceImpl.Create",
			senderErr,
		)
	}
	if sender == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	question, questionErr := tx.QuestionRepo().GetByID(dto.QuestionID, nil)
	if questionErr != nil {
		return nil, types.NewServerError(
			"Error in fetching question by id",
			"QuestionAnswerServiceImpl.Create",
			questionErr,
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
	if err := tx.AnswerRepo().Create(answer); err != nil {
		return nil, types.NewServerError(
			"Error in creating answer",
			"QuestionAnswerServiceImpl.Create",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return answer, nil
}

func (svc QuestionAnswerServiceImpl) GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	question, questionErr := tx.QuestionRepo().GetByID(questionID, nil)
	if questionErr != nil {
		return nil, types.NewServerError(
			"Error in fetching question by id",
			"QuestionAnswerServiceImpl.GetQuestionAnswer",
			questionErr,
		)
	}
	if question == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("question.errors.not_found"),
		)
	}
	answers, answersErr := tx.AnswerRepo().GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"question_id": questionID,
		},
		Relations: []string{"Sender"},
	})
	if answersErr != nil {
		return nil, types.NewServerError(
			"Error in fetching answers",
			"QuestionAnswerServiceImpl.GetQuestionAnswers",
			answersErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return answers, nil
}
