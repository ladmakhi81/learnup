package dtoreq

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
)

type AddVideoToCourseReq struct {
	CourseID    uint                      `json:"courseId" validate:"required,gte=1,numeric"`
	Title       string                    `json:"title" validate:"required,min=3"`
	Description string                    `json:"description" validate:"required,min=10"`
	AccessLevel entities.VideoAccessLevel `json:"accessLevel" validate:"required,oneof=private public"`
	IsPublished bool                      `json:"isPublished" validate:"required,boolean"`
}
