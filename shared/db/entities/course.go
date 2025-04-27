package entities

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model

	Name                        string                  `gorm:"column:name;index;not null;type:varchar(255)"`
	TeacherID                   *uint                   `gorm:"column:teacher_id;type:int unsigned;not null"`
	Teacher                     *User                   `gorm:"foreignKey:teacher_id;"`
	CategoryID                  *uint                   `gorm:"column:category_id;type:int unsigned; not null;"`
	Category                    *Category               `gorm:"foreignKey:category_id"`
	Price                       float64                 `gorm:"column:price;type:decimal(10,2);not null"`
	ThumbnailImage              string                  `gorm:"column:thumbnail_image;type:text;not null"`
	Image                       string                  `gorm:"column:image;type:text;not null"`
	Description                 string                  `gorm:"column:description;type:text;not null"`
	Prerequisite                string                  `gorm:"column:prerequisite;type:text;not null"`
	Level                       CourseLevel             `gorm:"column:level;type:text;not null"`
	Status                      CourseStatus            `gorm:"column:status;type:text;not null;default:'starting'"`
	StatusChangedAt             *time.Time              `gorm:"column:status_changed_at;type:timestamp;"`
	Tags                        []string                `gorm:"column:tags;type:text;serializer:json"`
	AbilityToAddComment         bool                    `gorm:"column:ability_to_add_comment;type:boolean;default:false"`
	CommentAccessMode           CourseCommentAccessMode `gorm:"column:comment_access_mode;type:text;not null;default:'all'"`
	IsPublished                 bool                    `gorm:"column:is_published;type:boolean;not null;default:false"`
	IsVerifiedByAdmin           bool                    `gorm:"column:is_verified_by_admin;type:boolean;not null;default:false"`
	VerifiedByID                *uint                   `gorm:"column:verified_by_id;type:int unsigned;"`
	VerifiedBy                  *User                   `gorm:"foreignKey:verified_by_id;"`
	VerifiedDate                *time.Time              `gorm:"column:verified_date;type:timestamp;"`
	Fee                         float64                 `gorm:"column:fee;type:decimal(10,2);not null;default:0"`
	IntroductionVideo           string                  `gorm:"column:introduction_video;type:text;not null"`
	CanHaveDiscount             bool                    `gorm:"column:can_have_discount;type:boolean;not null;default:false"`
	MaxDiscountAmount           float64                 `gorm:"column:max_discount_amount;type:decimal(10,2);not null;default:0"`
	DiscountFeeAmountPercentage float64                 `gorm:"column:discount_fee_amount_percentage;type:float;not null;default:0"`
	Participants                []*CourseParticipant    `gorm:"foreignKey:course_id"`
	ForumID                     *uint                   `gorm:"column:forum_id;type:int;"`
	Forum                       *CourseForum            `gorm:"foreignKey:forum_id"`
}

func (Course) TableName() string {
	return "_courses"
}

func (course Course) CalculateTeacherIncome() float64 {
	return course.Price - course.Fee
}

func (course *Course) SetPrice(price *float64) {
	if price == nil {
		course.Price = 0
	} else {
		course.Price = *price
	}
}

func (course *Course) SetFee(fee *float64) {
	if fee == nil {
		course.Fee = 0
	} else {
		course.Fee = *fee
	}
}

func (course Course) CheckFee(fee float64) bool {
	return fee > course.Price || fee < 0 || fee > course.Price-course.MaxDiscountAmount
}

func (course Course) IsTeacher(userID uint) bool {
	return *course.TeacherID == userID
}
