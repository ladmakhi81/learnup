package db

import (
	categoryEntity "github.com/ladmakhi81/learnup/internals/category/entity"
	userEntity "github.com/ladmakhi81/learnup/internals/user/entity"
)

func LoadEntities() map[string]any {
	return map[string]any{
		"user":     &userEntity.User{},
		"category": &categoryEntity.Category{},
	}
}
