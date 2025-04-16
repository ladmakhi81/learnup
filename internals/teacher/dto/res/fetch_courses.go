package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
	"time"
)

type FetchCourseItemRes struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
	Name        string
	Description string
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
