package db

import "github.com/ladmakhi81/learnup/internals/user/entity"

func LoadEntities() map[string]any {
	return map[string]any{
		"user": &entity.User{},
	}
}
