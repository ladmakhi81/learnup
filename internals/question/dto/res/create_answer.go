package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"time"
)

type CreateAnswerRes struct {
	ID         uint      `json:"id"`
	CreateAt   time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Content    string    `json:"content"`
	QuestionID uint      `json:"questionId"`
	SenderID   uint      `json:"senderId"`
}

func NewCreateAnswerRes(answer *entities.QuestionAnswer) CreateAnswerRes {
	return CreateAnswerRes{
		ID:         answer.ID,
		CreateAt:   answer.CreatedAt,
		UpdatedAt:  answer.UpdatedAt,
		Content:    answer.Content,
		QuestionID: answer.QuestionID,
		SenderID:   answer.SenderID,
	}
}
