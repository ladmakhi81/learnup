package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
	"time"
)

type getUserItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type getCourseItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetCommentPageItemDto struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	Content   string         `json:"content"`
	User      *getUserItem   `json:"user"`
	Course    *getCourseItem `json:"course"`
	IsReport  bool           `json:"isReport"`
}

func MapGetCommentPageItemsDto(comments []*entities.Comment) []*GetCommentPageItemDto {
	result := make([]*GetCommentPageItemDto, len(comments))
	for index, comment := range comments {
		result[index] = &GetCommentPageItemDto{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			DeletedAt: comment.DeletedAt,
			Content:   comment.Content,
			User: &getUserItem{
				ID:       comment.User.ID,
				Phone:    comment.User.Phone,
				FullName: comment.User.FullName(),
			},
			Course: &getCourseItem{
				ID:          comment.Course.ID,
				Name:        comment.Course.Name,
				Description: comment.Course.Description,
			},
			IsReport: comment.IsReport,
		}
	}
	return result
}
