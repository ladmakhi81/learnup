package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
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

type GetCommentPageableItemRes struct {
	ID        uint                  `json:"id"`
	Content   string                `json:"content"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	User      *getCommentUserItem   `json:"user"`
	Course    *getCommentCourseItem `json:"course"`
}

func MappedGetCommentPageableItemsRes(comments []*entities.Comment) []*GetCommentPageableItemRes {
	res := make([]*GetCommentPageableItemRes, len(comments))
	for i, comment := range comments {
		res[i] = &GetCommentPageableItemRes{
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
