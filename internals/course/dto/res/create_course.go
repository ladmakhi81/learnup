package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"time"
)

type CreateCourseResDto struct {
	ID                          uint                             `json:"id"`
	Name                        string                           `json:"name"`
	Description                 string                           `json:"description"`
	TeacherID                   *uint                            `json:"teacherId"`
	CategoryID                  *uint                            `json:"category_id"`
	Price                       float64                          `json:"price"`
	ThumbnailImage              string                           `json:"thumbnailImage"`
	Image                       string                           `json:"image"`
	Prerequisite        string                           `json:"prerequisite"`
	Level               entities.CourseLevel             `json:"level"`
	Status              entities.CourseStatus            `json:"courseStatus"`
	Tags                []string                         `json:"tags"`
	AbilityToAddComment bool                             `json:"abilityToAddComment"`
	CommentAccessMode   entities.CourseCommentAccessMode `json:"commentAccessMode"`
	IsPublished         bool                             `json:"isPublished"`
	IsVerifiedByAdmin           bool                             `json:"isVerifiedByAdmin"`
	VerifiedByID                *uint                            `json:"verifiedByID"`
	VerifiedDate                *time.Time                       `json:"verifiedDate"`
	Fee                         float64                          `json:"fee"`
	IntroductionVideo           string                           `json:"introductionVideo"`
	CanHaveDiscount             bool                             `json:"canHaveDiscount"`
	MaxDiscountAmount           float64                          `json:"maxDiscountAmount"`
	DiscountFeeAmountPercentage float64                          `json:"discountFeeAmountPercentage"`
	TeacherIncomeAmount         float64                          `json:"teacherIncomeAmount"`
}

func NewCreateCourseResDto(course *entities.Course) CreateCourseResDto {
	return CreateCourseResDto{
		ID:                          course.ID,
		Fee:                         course.Fee,
		Price:                       course.Price,
		VerifiedByID:                course.VerifiedByID,
		VerifiedDate:                course.VerifiedDate,
		TeacherID:                   course.TeacherID,
		ThumbnailImage:              course.ThumbnailImage,
		Tags:                        course.Tags,
		Status:                      course.Status,
		Prerequisite:                course.Prerequisite,
		MaxDiscountAmount:           course.MaxDiscountAmount,
		Level:                       course.Level,
		IsVerifiedByAdmin:           course.IsVerifiedByAdmin,
		IntroductionVideo:           course.IntroductionVideo,
		Image:                       course.Image,
		IsPublished:                 course.IsPublished,
		DiscountFeeAmountPercentage: course.DiscountFeeAmountPercentage,
		Description:                 course.Description,
		CommentAccessMode:           course.CommentAccessMode,
		CanHaveDiscount:             course.CanHaveDiscount,
		AbilityToAddComment:         course.AbilityToAddComment,
		Name:                        course.Name,
		CategoryID:                  course.CategoryID,
	}
}
