package entities

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	gorm.Model
	Content         string     `gorm:"column:content;type:text;not null;"`
	UserID          *uint      `gorm:"column:user_id;not null;index;type:int;"`
	User            *User      `gorm:"foreignkey:user_id;"`
	CourseID        *uint      `gorm:"column:course_id;not null;index;type:int;"`
	Course          *Course    `gorm:"foreignkey:course_id;"`
	IsReport        bool       `gorm:"column:is_report;type:boolean;default:false;"`
	ReportByID      *uint      `gorm:"column:report_by_id;type:int;index;"`
	ReportBy        *User      `gorm:"foreignkey:report_by_id;"`
	ReportDate      *time.Time `gorm:"column:report_date;type:timestamp;;"`
	ParentCommentId *uint      `gorm:"column:parent_comment_id;type:int;index;"`
	ParentComment   *Comment   `gorm:"foreignkey:parent_comment_id;"`
	Children        []*Comment `gorm:"foreignkey:parent_comment_id;"`
}

func (Comment) TableName() string {
	return "_comments"
}
