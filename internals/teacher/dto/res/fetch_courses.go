package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
	"time"
)

type FetchCourseItemDto struct {
	ID          uint           `json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

func MapFetchCourseItemsDto(courses []*entities.Course) []*FetchCourseItemDto {
	res := make([]*FetchCourseItemDto, len(courses))
	for index, course := range courses {
		res[index] = &FetchCourseItemDto{
			ID:          course.ID,
			CreatedAt:   course.CreatedAt,
			UpdatedAt:   course.UpdatedAt,
			DeletedAt:   course.DeletedAt,
			Name:        course.Name,
			Description: course.Description,
		}
	}
	return res
}
