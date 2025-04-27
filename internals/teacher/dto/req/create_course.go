package dtoreq

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
)

type CreateCourseReqDto struct {
	Name                string                           `json:"name" validate:"required,min=3,max=255"`
	CategoryID          uint                             `json:"categoryId" validate:"required,numeric"`
	Price               float64                          `json:"price,omitempty" validate:"gte=0"`
	ThumbnailImage      string                           `json:"thumbnailImage" validate:"required,min=10"`
	Image               string                           `json:"image" validate:"required,min=10"`
	Description         string                           `json:"description" validate:"required,min=20"`
	Prerequisite        string                           `json:"prerequisite" validate:"required,min=20"`
	Level               entities.CourseLevel             `json:"courseLevel" validate:"required,oneof=beginner pre-intermediate intermediate advance"`
	Tags                []string                         `json:"tags,omitempty"`
	AbilityToAddComment bool                             `json:"abilityToAddComment" validate:"required"`
	CommentAccessMode   entities.CourseCommentAccessMode `json:"commentAccessMode,omitempty" validate:"oneof=all students"`
	CanHaveDiscount     bool                             `json:"canHaveDiscount" validate:"required"`
	MaxDiscountAmount   float64                          `json:"maxDiscountAmount,omitempty" validate:"numeric,gte=0"`
	IntroductionVideo   string                           `json:"introductionVideo" validate:"required,min=20"`
	ForumStatus         entities.CourseForumStatus       `json:"forumStatus" validate:"required,oneof=open close"`
	IsForumPublic       bool                             `json:"isForumPublic" validate:"required,boolean"`
	ForumAccessMode     entities.CourseForumAccessMode   `json:"forumAccessMode" validate:"required,oneof=student-only teacher-only student-teacher"`
}
