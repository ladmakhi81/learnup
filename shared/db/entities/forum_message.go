package entities

import "gorm.io/gorm"

type ForumMessage struct {
	gorm.Model
	Message string      `gorm:"column:message;type:text;not null"`
	UserID  uint        `gorm:"column:user_id;type:int;index;not null"`
	User    *User       `gorm:"foreignKey:user_id"`
	ForumID uint        `gorm:"column:forum_id;type:int;index;not null"`
	Forum   CourseForum `gorm:"foreignKey:forum_id"`
}

func (ForumMessage) TableName() string {
	return "_forum_messages"
}
