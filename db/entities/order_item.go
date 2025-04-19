package entities

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID  uint    `gorm:"column:order_id;type:int;not null;index"`
	Order    *Order  `gorm:"foreignKey:order_id;"`
	UserID   uint    `gorm:"column:user_id;type:int;not null;index"`
	CourseID uint    `gorm:"column:course_id;type:int;not null;index"`
	Course   *Course `gorm:"foreignKey:course_id;"`
	Amount   float64 `gorm:"column:amount;type:decimal(10,2);"`
}

func (OrderItem) TableName() string {
	return "_orders"
}
