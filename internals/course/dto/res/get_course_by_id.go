package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type teacherUser struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
}

type categoryItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	IsPublished bool   `json:"isPublished"`
}

type verifiedByItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
}

type GetCourseByIdRes struct {
	ID                          uint                              `json:"id"`
	CreatedAt                   time.Time                         `json:"createdAt"`
	UpdatedAt                   time.Time                         `json:"updatedAt"`
	DeletedAt                   time.Time                         `json:"deletedAt"`
	Teacher                     *teacherUser                      `json:"teacher"`
	Category                    *categoryItem                     `json:"category"`
	Price                       float64                           `json:"price"`
	ThumbnailImage              string                            `json:"thumbnailImage"`
	Image                       string                            `json:"image"`
	Description                 string                            `json:"description"`
	Prerequisite                string                            `json:"prerequisite"`
	Level                       entities2.CourseLevel             `json:"level"`
	Status                      entities2.CourseStatus            `json:"status"`
	StatusChangedAt             *time.Time                        `json:"statusChangedAt"`
	Tags                        []string                          `json:"tags"`
	AbilityToAddComment         bool                              `json:"abilityToAddComment"`
	CommentAccessMode           entities2.CourseCommentAccessMode `json:"commentAccessMode"`
	IsPublished                 bool                              `json:"isPublished"`
	IsVerifiedByAdmin           bool                              `json:"isVerifiedByAdmin"`
	VerifiedBy                  *verifiedByItem                   `json:"verifiedBy"`
	VerifiedDate                *time.Time                        `json:"verifiedDate"`
	Fee                         float64                           `json:"fee"`
	IntroductionVideo           string                            `json:"introductionVideo"`
	CanHaveDiscount             bool                              `json:"canHaveDiscount"`
	MaxDiscountAmount           float64                           `json:"maxDiscountAmount"`
	DiscountFeeAmountPercentage float64                           `json:"discountFeeAmountPercentage"`
}

func NewGetCourseByIdRes(course *entities2.Course) GetCourseByIdRes {
	res := GetCourseByIdRes{
		ID:                course.ID,
		Status:            course.Status,
		VerifiedDate:      course.VerifiedDate,
		CreatedAt:         course.CreatedAt,
		UpdatedAt:         course.UpdatedAt,
		IsPublished:       course.IsPublished,
		DeletedAt:         course.DeletedAt.Time,
		Description:       course.Description,
		StatusChangedAt:   course.StatusChangedAt,
		IsVerifiedByAdmin: course.IsVerifiedByAdmin,
		IntroductionVideo: course.IntroductionVideo,
		Image:             course.Image,
		ThumbnailImage:    course.ThumbnailImage,
		Price:             course.Price,
		Category: &categoryItem{
			ID:          course.Category.ID,
			Name:        course.Category.Name,
			IsPublished: course.Category.IsPublished,
		},
		Teacher: &teacherUser{
			ID:       course.Teacher.ID,
			FullName: course.Teacher.FullName(),
		},
		AbilityToAddComment:         course.AbilityToAddComment,
		CanHaveDiscount:             course.CanHaveDiscount,
		CommentAccessMode:           course.CommentAccessMode,
		DiscountFeeAmountPercentage: course.DiscountFeeAmountPercentage,
		Fee:                         course.Fee,
		Level:                       course.Level,
		MaxDiscountAmount:           course.MaxDiscountAmount,
		Prerequisite:                course.Prerequisite,
		Tags:                        course.Tags,
	}

	if course.VerifiedBy != nil {
		res.VerifiedBy = &verifiedByItem{
			ID:       course.VerifiedBy.ID,
			FullName: course.VerifiedBy.FullName(),
		}
	}

	return res
}
