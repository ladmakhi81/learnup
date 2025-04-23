package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
)

type CreateCommentResDto struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

func NewCreateCommentResDto(comment *entities.Comment) *CreateCommentResDto {
	return &CreateCommentResDto{
		ID:      comment.ID,
		Content: comment.Content,
	}
}
