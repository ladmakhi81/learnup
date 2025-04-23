package dtores

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type getCommentUserItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type getCommentCourseItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetCommentPageableItemDto struct {
	ID        uint                  `json:"id"`
	Content   string                `json:"content"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	User      *getCommentUserItem   `json:"user"`
	Course    *getCommentCourseItem `json:"course"`
}

func MapGetCommentPageableItemsDto(comments []*entities.Comment) []*GetCommentPageableItemDto {
	res := make([]*GetCommentPageableItemDto, len(comments))
	for i, comment := range comments {
		res[i] = &GetCommentPageableItemDto{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			User: &getCommentUserItem{
				ID:       comment.User.ID,
				FullName: comment.User.FullName(),
				Phone:    comment.User.Phone,
			},
			Course: &getCommentCourseItem{
				ID:          comment.Course.ID,
				Name:        comment.Course.Name,
				Description: comment.Course.Description,
			},
		}
	}
	return res
}
