package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"time"
)

type CreateQuestionRes struct {
	ID        uint                      `json:"id"`
	CreatedAt time.Time                 `json:"createdAt"`
	Content   string                    `json:"content"`
	Priority  entities.QuestionPriority `json:"priority"`
}

func NewCreateQuestionRes(question *entities.Question) *CreateQuestionRes {
	res := &CreateQuestionRes{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		Content:   question.Content,
		Priority:  question.Priority,
	}
	return res
}
