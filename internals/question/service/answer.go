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
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewQuestionAnswerServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *QuestionAnswerServiceImpl {
	return &QuestionAnswerServiceImpl{
		translationSvc: translationSvc,
		repo:           repo,
	}
}

func (svc QuestionAnswerServiceImpl) Create(dto questionDtoReq.AnswerQuestionReq) (*entities.QuestionAnswer, error) {
	sender, senderErr := svc.repo.UserRepo.GetByID(dto.SenderID)
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
	question, questionErr := svc.repo.QuestionRepo.GetByID(dto.QuestionID)
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
	if err := svc.repo.AnswerRepo.Create(answer); err != nil {
		return nil, types.NewServerError(
			"Error in creating answer",
			"QuestionAnswerServiceImpl.Create",
			err,
		)
	}
	return answer, nil
}

func (svc QuestionAnswerServiceImpl) GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error) {
	question, questionErr := svc.repo.QuestionRepo.GetByID(questionID)
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
	answers, answersErr := svc.repo.AnswerRepo.GetAll(repositories.GetAllOptions{
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
	return answers, nil
}
