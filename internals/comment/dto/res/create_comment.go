package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
)

type CreateCommentRes struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

func NewCreateCommentRes(comment *entities.Comment) *CreateCommentRes {
	return &CreateCommentRes{
		ID:      comment.ID,
		Content: comment.Content,
	}
}
