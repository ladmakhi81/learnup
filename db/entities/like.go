package entities

import (
	"time"
)

type Like struct {
	ID        uint      `gorm:"column:id;primary_key"`
	UserID    uint      `gorm:"column:user_id;not null;index;type:int;"`
	User      *User     `gorm:"foreignkey:user_id"`
	CourseID  uint      `gorm:"column:course_id;not null;index;type:int;"`
	Course    *Course   `gorm:"foreignkey:course_id"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	Type      LikeType  `gorm:"column:type;not null;type:varchar(255);"`
}

func (Like) TableName() string {
	return "_likes"
}
