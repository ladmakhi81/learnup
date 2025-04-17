package db

import (
	"github.com/ladmakhi81/learnup/db/entities"
)

func LoadEntities() map[string]any {
	return map[string]any{
		"user":         &entities.User{},
		"category":     &entities.Category{},
		"course":       &entities.Course{},
		"video":        &entities.Video{},
		"notification": &entities.Notification{},
		"comment":      &entities.Comment{},
		"like":         &entities.Like{},
	}
}
