package db

import (
	entities2 "github.com/ladmakhi81/learnup/shared/db/entities"
)

func LoadEntities() map[string]any {
	return map[string]any{
		"user":            &entities2.User{},
		"category":        &entities2.Category{},
		"course":          &entities2.Course{},
		"video":           &entities2.Video{},
		"notification":    &entities2.Notification{},
		"comment":         &entities2.Comment{},
		"like":            &entities2.Like{},
		"question":        &entities2.Question{},
		"question_answer": &entities2.QuestionAnswer{},
		"cart":            &entities2.Cart{},
		"order":           &entities2.Order{},
		"order_items":     &entities2.OrderItem{},
		"payment":         &entities2.Payment{},
		"transaction":     &entities2.Transaction{},
	}
}
