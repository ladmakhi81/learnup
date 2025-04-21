package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type CreateQuestionRes struct {
	ID        uint                       `json:"id"`
	CreatedAt time.Time                  `json:"createdAt"`
	Content   string                     `json:"content"`
	Priority  entities2.QuestionPriority `json:"priority"`
}

func NewCreateQuestionRes(question *entities2.Question) *CreateQuestionRes {
	res := &CreateQuestionRes{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		Content:   question.Content,
		Priority:  question.Priority,
	}
	return res
}
