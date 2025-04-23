package dtores

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type CreateCourseResDto struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewCreateCourseResDto(course *entities.Course) CreateCourseResDto {
	return CreateCourseResDto{
		ID:        course.ID,
		CreatedAt: course.CreatedAt,
		UpdatedAt: course.UpdatedAt,
	}
}
