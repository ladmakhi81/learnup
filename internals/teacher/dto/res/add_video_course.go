package dtores

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
)

type AddVideoToCourseResDto struct {
	ID          uint                      `json:"id"`
	CourseID    uint                      `json:"courseId"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	AccessLevel entities.VideoAccessLevel `json:"accessLevel"`
	IsPublished bool                      `json:"isPublished"`
}

func NewAddVideoToCourseResDto(video *entities.Video) *AddVideoToCourseResDto {
	return &AddVideoToCourseResDto{
		ID:          video.ID,
		CourseID:    *video.CourseId,
		Title:       video.Title,
		Description: video.Description,
		AccessLevel: video.AccessLevel,
		IsPublished: video.IsPublished,
	}
}
