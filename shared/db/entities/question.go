package entities

import (
	"gorm.io/gorm"
	"time"
)

type Question struct {
	gorm.Model
	UserID   uint              `gorm:"column:user_id;not null;index;type:int;"`
	User     *User             `gorm:"foreignkey:user_id"`
	CourseID uint              `gorm:"column:course_id;not null;index;type:int;"`
	Course   *Course           `gorm:"foreignkey:course_id"`
	Content  string            `gorm:"column:content;type:text;not null"`
	Priority QuestionPriority  `gorm:"column:priority;type:varchar(255);not null"`
	IsClosed bool              `gorm:"column:is_closed;type:boolean;default:false"`
	ClosedDate *time.Time        `gorm:"column:closed_date;type:timestamp;default:null"`
	VideoID  *uint             `gorm:"column:video_id;index;type:int;default:null"`
	Video    *Video            `gorm:"foreignkey:video_id"`
	Answers  []*QuestionAnswer `gorm:"foreignkey:question_id"`
}

func (Question) TableName() string {
	return "_questions"
}
