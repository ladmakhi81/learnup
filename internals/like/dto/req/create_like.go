package dtoreq

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
)

type CreateLikeReqDto struct {
	Type     entities.LikeType `json:"type" validate:"required,oneof=none like dislike"`
	CourseID uint              `json:"-"`
}
