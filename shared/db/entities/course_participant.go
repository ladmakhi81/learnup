package entities

import (
	"time"
)

type CourseParticipant struct {
	CourseID           uint       `gorm:"column:course_id;type:int;index;"`
	StudentID          uint       `gorm:"column:student_id;type:int;index;"`
	Student            *User      `gorm:"foreignKey:student_id"`
	TeacherID          uint       `gorm:"column:teacher_id;type:int;not null"`
	LastVideoWatchDate *time.Time `gorm:"column:last_video_watch_date;type:timestamp;default:null"`
	CreatedAt          time.Time
}

func (CourseParticipant) TableName() string {
	return "_course_participants"
}
