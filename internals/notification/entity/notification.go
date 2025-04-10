package entity

import (
	userEntity "github.com/ladmakhi81/learnup/internals/user/entity"
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	gorm.Model

	Type     NotificationType `gorm:"column:type;type:varchar(255);not null;"`
	IsSeen   bool             `gorm:"column:is_seen;type:boolean;default:false;"`
	SeenAt   *time.Time       `gorm:"column:seen_at;type:timestamp;"`
	UserID   *uint            `gorm:"column:user_id;type:int;index;not null;"`
	User     *userEntity.User `gorm:"foreignKey:user_id;"`
	Metadata any              `gorm:"column:metadata;type:text;serializer:json;not null;"`
}

func (Notification) TableName() string {
	return "_notifications"
}
