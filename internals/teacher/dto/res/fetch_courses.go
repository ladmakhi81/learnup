package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
	"time"
)

type FetchCourseItemRes struct {
	ID          uint           `json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

func MapCoursesToFetchCourseItemRes(courses []*entities.Course) []*FetchCourseItemRes {
	res := make([]*FetchCourseItemRes, len(courses))
	for index, course := range courses {
		res[index] = &FetchCourseItemRes{
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
