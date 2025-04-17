package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
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

type GetQuestionItemRes struct {
	ID        uint                      `json:"id"`
	CreatedAt time.Time                 `json:"createdAt"`
	UpdatedAt time.Time                 `json:"updatedAt"`
	DeletedAt gorm.DeletedAt            `json:"deletedAt"`
	Content   string                    `json:"content"`
	Priority  entities.QuestionPriority `json:"priority"`
	User      *questionUserItem         `json:"user"`
	Course    *questionCourseItem       `json:"course"`
	Video     *questionVideoItem        `json:"video"`
}

func MapGetQuestionItemRes(questions []*entities.Question) []*GetQuestionItemRes {
	res := make([]*GetQuestionItemRes, len(questions))
	for index, question := range questions {
		res[index] = &GetQuestionItemRes{
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
