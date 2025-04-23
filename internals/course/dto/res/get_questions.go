package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
	"time"
)

type questionUserItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
}

type questionCourseItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type questionVideoItem struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetQuestionItemDto struct {
	ID        uint                       `json:"id"`
	CreatedAt time.Time                  `json:"createdAt"`
	UpdatedAt time.Time                  `json:"updatedAt"`
	DeletedAt gorm.DeletedAt             `json:"deletedAt"`
	Content  string                     `json:"content"`
	Priority entities2.QuestionPriority `json:"priority"`
	User     *questionUserItem          `json:"user"`
	Course    *questionCourseItem        `json:"course"`
	Video     *questionVideoItem         `json:"video"`
}

func MapGetQuestionItemsDto(questions []*entities2.Question) []*GetQuestionItemDto {
	res := make([]*GetQuestionItemDto, len(questions))
	for index, question := range questions {
		res[index] = &GetQuestionItemDto{
			ID:        question.ID,
			CreatedAt: question.CreatedAt,
			UpdatedAt: question.UpdatedAt,
			DeletedAt: question.DeletedAt,
			Content:   question.Content,
			Priority:  question.Priority,
			Course: &questionCourseItem{
				ID:          question.Course.ID,
				Name:        question.Course.Name,
				Description: question.Course.Description,
			},
			User: &questionUserItem{
				ID:       question.User.ID,
				FullName: question.User.FullName(),
			},
		}
		if question.Video != nil {
			res[index].Video = &questionVideoItem{
				ID:          question.Video.ID,
				Title:       question.Video.Title,
				Description: question.Video.Description,
			}
		}
	}
	return res
}
