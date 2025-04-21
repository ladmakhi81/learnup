package entities

import "gorm.io/gorm"

type QuestionAnswer struct {
	gorm.Model
	SenderID   uint      `gorm:"type:int;not null;index;column:sender_id;"`
	Sender     *User     `gorm:"foreignkey:sender_id"`
	Content    string    `gorm:"type:text;not null;column:content"`
	QuestionID uint      `gorm:"type:int;not null;index;column:question_id"`
	Question   *Question `gorm:"foreignkey:question_id"`
}

func (QuestionAnswer) TableName() string {
	return "_question_answers"
}
