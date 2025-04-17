package dtoreq

import "github.com/ladmakhi81/learnup/db/entities"

type CreateLikeReq struct {
	Type     entities.LikeType `json:"type" validate:"required,oneof=none like dislike"`
	CourseID uint              `json:"-"`
}
