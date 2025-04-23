package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"time"
)

type CreateQuestionResDto struct {
	ID        uint                      `json:"id"`
	CreatedAt time.Time                 `json:"createdAt"`
	Content  string                    `json:"content"`
	Priority entities.QuestionPriority `json:"priority"`
}

func NewCreateQuestionResDto(question *entities.Question) CreateQuestionResDto {
	return CreateQuestionResDto{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		Content:   question.Content,
		Priority:  question.Priority,
	}
}
