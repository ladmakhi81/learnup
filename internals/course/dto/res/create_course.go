package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"time"
)

type CreateCourseRes struct {
	ID                          uint                             `json:"id"`
	Name                        string                           `json:"name"`
	Description                 string                           `json:"description"`
	TeacherID                   *uint                            `json:"teacherId"`
	CategoryID                  *uint                            `json:"category_id"`
	Price                       float64                          `json:"price"`
	ThumbnailImage              string                           `json:"thumbnailImage"`
	Image                       string                           `json:"image"`
	Prerequisite                string                           `json:"prerequisite"`
	Level                       entities.CourseLevel             `json:"level"`
	Status                      entities.CourseStatus            `json:"courseStatus"`
	Tags                        []string                         `json:"tags"`
	AbilityToAddComment         bool                             `json:"abilityToAddComment"`
	CommentAccessMode           entities.CourseCommentAccessMode `json:"commentAccessMode"`
	IsPublished                 bool                             `json:"isPublished"`
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
