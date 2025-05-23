package db

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
)

func LoadEntities() map[string]any {
	return map[string]any{
		"user":               &entities.User{},
		"category":           &entities.Category{},
		"course":             &entities.Course{},
		"video":              &entities.Video{},
		"notification":       &entities.Notification{},
		"comment":            &entities.Comment{},
		"like":               &entities.Like{},
		"question":           &entities.Question{},
		"question_answer":    &entities.QuestionAnswer{},
		"cart":               &entities.Cart{},
		"order":              &entities.Order{},
		"order_items":        &entities.OrderItem{},
		"payment":            &entities.Payment{},
		"transaction":        &entities.Transaction{},
		"course_forum":       &entities.CourseForum{},
		"course_participant": &entities.CourseParticipant{},
		"course_message":     &entities.ForumMessage{},
	}
}
