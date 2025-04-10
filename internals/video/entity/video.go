package entity

import (
	courseEntity "github.com/ladmakhi81/learnup/internals/course/entity"
	userEntity "github.com/ladmakhi81/learnup/internals/user/entity"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model

	CourseId     *uint                `gorm:"column:course_id;type:int;index;"`
	Course       *courseEntity.Course `gorm:"foreignKey:course_id"`
	Title        string               `gorm:"column:title;type:varchar(255);unique;not null;"`
	Description  string               `gorm:"column:description;type:text;not null;"`
	AccessLevel  VideoAccessLevel     `gorm:"column:access_level;type:varchar(255);not null;"`
	IsPublished  bool                 `gorm:"column:is_published;type:boolean;default:false;"`
	IsVerified   bool                 `gorm:"column:is_verified;type:boolean;default:false;"`
	VerifiedDate *time.Time           `gorm:"column:verified_date;type:timestamp;"`
	VerifiedById *uint                `gorm:"column:verified_by_id;type:int;index;"`
	VerifiedBy   *userEntity.User     `gorm:"foreignKey:verified_by_id;"`
	Duration     *float64             `gorm:"column:duration;type:decimal(10,2);"`
	Status       VideoStatus          `gorm:"column:status;type:varchar(255);"`
	URL          string               `gorm:"column:video_url;type:text;"`
}

func (Video) TableName() string {
	return "_videos"
}
