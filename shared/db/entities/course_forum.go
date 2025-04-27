package entities

import (
	"gorm.io/gorm"
	"time"
)

type CourseForum struct {
	gorm.Model
	CourseID        uint                  `gorm:"column:course_id;type:int;not null;"`
	Course          *Course               `gorm:"foreignKey:course_id;"`
	TeacherID       uint                  `gorm:"column:teacher_id;type:int;not null;"`
	Teacher         *User                 `gorm:"foreignKey:teacher_id"`
	Status          CourseForumStatus     `gorm:"column:status;type:varchar(255);default:'open';"`
	StatusChangedAt *time.Time            `gorm:"column:status_changed_at;type:timestamp;default:null;"`
	IsPublic        bool                  `gorm:"column:is_public;type:boolean;default:false;"`
	AccessMode      CourseForumAccessMode `gorm:"column:access_mode;type:varchar(255);default:'teacher-only';"`
	Messages        []*ForumMessage       `gorm:"foreignKey:forum_id"`
}

func (CourseForum) TableName() string {
	return "_course_forums"
}
