package entities

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model

	CourseId     *uint            `gorm:"column:course_id;type:int;index;"`
	Course       *Course          `gorm:"foreignKey:course_id"`
	Title        string           `gorm:"column:title;type:varchar(255);index;not null;"`
	Description  string           `gorm:"column:description;type:text;not null;"`
	AccessLevel  VideoAccessLevel `gorm:"column:access_level;type:varchar(255);not null;"`
	IsPublished  bool             `gorm:"column:is_published;type:boolean;default:false;"`
	IsVerified   bool             `gorm:"column:is_verified;type:boolean;default:false;"`
	VerifiedDate *time.Time       `gorm:"column:verified_date;type:timestamp;"`
	VerifiedById *uint            `gorm:"column:verified_by_id;type:int;index;"`
	VerifiedBy   *User            `gorm:"foreignKey:verified_by_id;"`
	Duration     *string          `gorm:"column:duration;type:text;"`
	Status       VideoStatus      `gorm:"column:status;type:varchar(255);"`
	URL          string           `gorm:"column:video_url;type:text;"`
}

func (Video) TableName() string {
	return "_videos"
}
