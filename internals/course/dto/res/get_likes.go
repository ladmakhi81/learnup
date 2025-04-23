package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
)

type userItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
}

type GetLikesPageableItemDto struct {
	User *userItem         `json:"user"`
	Type entities.LikeType `json:"type"`
}

func MapGetLikesPageableItemsDto(likes []*entities.Like) []*GetLikesPageableItemDto {
	res := make([]*GetLikesPageableItemDto, len(likes))
	for i, like := range likes {
		res[i] = &GetLikesPageableItemDto{
			User: &userItem{
				ID:       like.UserID,
				FullName: like.User.FullName(),
			},
			Type: like.Type,
		}
	}
	return res
}
