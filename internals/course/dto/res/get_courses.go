package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type courseTeacher struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName'"`
	Phone    string `json:"phone"`
}

type courseUserVerifier struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName'"`
	Phone    string `json:"phone"`
}

type courseCategory struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	IsPublished bool   `json:"isPublished"`
}

type courseItems struct {
	ID                uint                   `json:"id"`
	Name              string                 `json:"name"`
	Teacher           *courseTeacher         `json:"teacher"`
	Category          *courseCategory        `json:"category"`
	Price             float64                `json:"price"`
	ThumbnailImage    string                 `json:"thumbnail"`
	IntroductionVideo string                 `json:"introductionVideo"`
	Status            entities2.CourseStatus `json:"status"`
	IsPublished       bool                   `json:"isPublished"`
	IsVerifiedByAdmin bool                   `json:"isVerified"`
	VerifiedDate      *time.Time             `json:"verifiedDate"`
	VerifiedBy        *courseUserVerifier    `json:"verifiedBy"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
	DeletedAt         time.Time              `json:"deletedAt"`
	StatusChangedAt   *time.Time             `json:"statusChangedAt"`
}

type GetCoursesRes struct {
	Courses []*courseItems `json:"courses"`
}

func mapCoursesToCoursesPageableItems(courses []*entities2.Course) []*courseItems {
	mappedCourses := make([]*courseItems, len(courses))
	for courseIndex, course := range courses {
		mappedCourses[courseIndex] = &courseItems{
			ID:   course.ID,
			Name: course.Name,
			Teacher: &courseTeacher{
				ID:       course.Teacher.ID,
				FullName: course.Teacher.FullName(),
				Phone:    course.Teacher.Phone,
			},
			Category: &courseCategory{
				ID:          course.Category.ID,
				Name:        course.Category.Name,
				IsPublished: course.Category.IsPublished,
			},
			Price:             course.Price,
			ThumbnailImage:    course.ThumbnailImage,
			IntroductionVideo: course.IntroductionVideo,
			Status:            course.Status,
			IsPublished:       course.IsPublished,
			IsVerifiedByAdmin: course.IsVerifiedByAdmin,
			VerifiedDate:      course.VerifiedDate,
			DeletedAt:         course.DeletedAt.Time,
			StatusChangedAt:   course.StatusChangedAt,
			CreatedAt:         course.CreatedAt,
			UpdatedAt:         course.UpdatedAt,
			VerifiedBy: &courseUserVerifier{
				ID:       course.VerifiedBy.ID,
				FullName: course.VerifiedBy.FullName(),
				Phone:    course.VerifiedBy.Phone,
			},
		}
	}
	return mappedCourses
}

func NewGetCoursesRes(courses []*entities2.Course) GetCoursesRes {
	res := GetCoursesRes{}
	res.Courses = mapCoursesToCoursesPageableItems(courses)
	return res
}
