package db

import (
	categoryEntity "github.com/ladmakhi81/learnup/internals/category/entity"
	courseEntity "github.com/ladmakhi81/learnup/internals/course/entity"
	userEntity "github.com/ladmakhi81/learnup/internals/user/entity"
	videoEntity "github.com/ladmakhi81/learnup/internals/video/entity"
)

func LoadEntities() map[string]any {
	return map[string]any{
		"user":     &userEntity.User{},
		"category": &categoryEntity.Category{},
		"course":   &courseEntity.Course{},
		"video":    &videoEntity.Video{},
	}
}
