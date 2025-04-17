package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	questionDtoReq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	questionRepository "github.com/ladmakhi81/learnup/internals/question/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type QuestionAnswerService interface {
	Create(dto questionDtoReq.AnswerQuestionReq) (*entities.QuestionAnswer, error)
	GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error)
}

type QuestionAnswerServiceImpl struct {
	answerRepo     questionRepository.QuestionAnswerRepo
	questionSvc    QuestionService
	userSvc        userService.UserSvc
	translationSvc contracts.Translator
}

func NewQuestionAnswerServiceImpl(
	answerRepo questionRepository.QuestionAnswerRepo,
	questionSvc QuestionService,
	translationSvc contracts.Translator,
	userSvc userService.UserSvc,
) *QuestionAnswerServiceImpl {
	return &QuestionAnswerServiceImpl{
		answerRepo:     answerRepo,
		questionSvc:    questionSvc,
		translationSvc: translationSvc,
		userSvc:        userSvc,
	}
}

func (svc QuestionAnswerServiceImpl) Create(dto questionDtoReq.AnswerQuestionReq) (*entities.QuestionAnswer, error) {
	sender, senderErr := svc.userSvc.FindById(dto.SenderID)
	if senderErr != nil {
		return nil, senderErr
	}
	if sender == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	question, questionErr := svc.questionSvc.FindById(dto.QuestionID)
	if questionErr != nil {
		return nil, questionErr
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
	if err := svc.answerRepo.Create(answer); err != nil {
		return nil, types.NewServerError(
			"Error in creating answer",
			"QuestionAnswerServiceImpl.Create",
			err,
		)
	}
	return answer, nil
}

func (svc QuestionAnswerServiceImpl) GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error) {
	question, questionErr := svc.questionSvc.FindById(questionID)
	if questionErr != nil {
		return nil, questionErr
	}
	if question == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("question.errors.not_found"),
		)
	}
	answers, answersErr := svc.answerRepo.Fetch(
		questionRepository.FetchAnswerOption{
			QuestionID: &question.ID,
			Preloads:   []string{"Sender"},
		},
	)
	if answersErr != nil {
		return nil, types.NewServerError(
			"Error in fetching answers",
			"QuestionAnswerServiceImpl.GetQuestionAnswers",
			answersErr,
		)
	}
	return answers, nil
}
