package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
)

type userItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
}

type GetLikesPageableItem struct {
	User *userItem          `json:"user"`
	Type entities2.LikeType `json:"type"`
}

func MappedGetLikesPageableItem(likes []*entities2.Like) []*GetLikesPageableItem {
	res := make([]*GetLikesPageableItem, len(likes))
	for i, like := range likes {
		res[i] = &GetLikesPageableItem{
			User: &userItem{
				ID:       like.UserID,
				FullName: like.User.FullName(),
			},
			Type: like.Type,
		}
	}
	return res
}
