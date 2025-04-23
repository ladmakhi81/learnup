package service

import (
	questionDtoReq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	questionError "github.com/ladmakhi81/learnup/internals/question/error"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/db/repositories"
	"github.com/ladmakhi81/learnup/shared/types"
)

type QuestionAnswerService interface {
	Create(sender *entities.User, dto questionDtoReq.AnswerQuestionReqDto) (*entities.QuestionAnswer, error)
	GetQuestionAnswers(questionID uint) ([]*entities.QuestionAnswer, error)
}

type questionAnswerService struct {
	unitOfWork db.UnitOfWork
}

func NewQuestionAnswerSvc(unitOfWork db.UnitOfWork) QuestionAnswerService {
	return &questionAnswerService{unitOfWork: unitOfWork}
}

func (svc questionAnswerService) Create(sender *entities.User, dto questionDtoReq.AnswerQuestionReqDto) (*entities.QuestionAnswer, error) {
	const operationName = "questionAnswerService.Create"
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
