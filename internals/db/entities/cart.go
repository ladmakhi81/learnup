package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID   uint    `gorm:"column:user_id;type:int;not null;index"`
	User     *User   `gorm:"foreignkey:user_id"`
	CourseID uint    `gorm:"column:course_id;type:int;not null;index"`
	Course   *Course `gorm:"foreignkey:course_id"`
}

func (Cart) TableName() string {
	return "_carts"
}

func (cart Cart) IsOwner(userID uint) bool {
	return cart.UserID == userID
}
