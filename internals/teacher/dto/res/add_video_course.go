package dtores

import "github.com/ladmakhi81/learnup/db/entities"

type AddVideoToCourseRes struct {
	ID          uint                      `json:"id"`
	CourseID    uint                      `json:"courseId"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	AccessLevel entities.VideoAccessLevel `json:"accessLevel"`
	IsPublished bool                      `json:"isPublished"`
}
