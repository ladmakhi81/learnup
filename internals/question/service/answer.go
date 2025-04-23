package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	questionDtoReq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	questionError "github.com/ladmakhi81/learnup/internals/question/error"
	userError "github.com/ladmakhi81/learnup/internals/user/error"
	"github.com/ladmakhi81/learnup/types"
)

type QuestionAnswerService interface {
	Create(dto questionDtoReq.AnswerQuestionReq) (*entities.QuestionAnswer, error)
	GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error)
}

type questionAnswerService struct {
	unitOfWork db.UnitOfWork
}

func NewQuestionAnswerSvc(unitOfWork db.UnitOfWork) QuestionAnswerService {
	return &questionAnswerService{unitOfWork: unitOfWork}
}

func (svc questionAnswerService) Create(dto questionDtoReq.AnswerQuestionReq) (*entities.QuestionAnswer, error) {
	const operationName = "questionAnswerService.Create"
	sender, err := svc.unitOfWork.UserRepo().GetByID(dto.SenderID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching user", operationName, err)
	}
	if sender == nil {
		return nil, userError.User_NotFound
	}
	question, err := svc.unitOfWork.QuestionRepo().GetByID(dto.QuestionID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching question by id", operationName, err)
	}
	if question == nil {
		return nil, questionError.Question_NotFound
	}
	if question.IsClosed {
		return nil, questionError.Question_ClosedStatus
	}
	answer := &entities.QuestionAnswer{
		QuestionID: dto.QuestionID,
		Content:    dto.Content,
		SenderID:   sender.ID,
	}
	if err := svc.unitOfWork.AnswerRepo().Create(answer); err != nil {
		return nil, types.NewServerError("Error in creating answer", operationName, err)
	}
	return answer, nil
}

func (svc questionAnswerService) GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error) {
	const operationName = "questionAnswerService.GetQuestionAnswers"
	question, err := svc.unitOfWork.QuestionRepo().GetByID(questionID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching question by id", operationName, err)
	}
	if question == nil {
		return nil, questionError.Question_NotFound
	}
	order := "created_at asc"
	answers, err := svc.unitOfWork.AnswerRepo().GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"question_id": questionID,
		},
		Relations: []string{"Sender"},
		Order:     &order,
	})
	if err != nil {
		return nil, types.NewServerError("Error in fetching answers", operationName, err)
	}
	return answers, nil
}
