package entities

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID          uint         `gorm:"column:user_id;type:int;not null;index"`
	User            *User        `gorm:"foreignkey:user_id"`
	TotalPrice      float64      `gorm:"type:decimal(10,2);default:0"`
	Status          OrderStatus  `gorm:"column:status;type:varchar(255);default:'pending';"`
	StatusChangedAt *time.Time   `gorm:"column:status_changed_at;type:timestamp;default:null"`
	FinalPrice      float64      `gorm:"type:decimal(10,2);default:0"`
	Items           []*OrderItem `gorm:"foreignkey:order_id"`
	DiscountPrice   float64      `gorm:"type:decimal(10,2);default:0"`
}

func (Order) TableName() string {
	return "_orders"
}

func NewOrder(userID uint) *Order {
	return &Order{
		UserID:        userID,
		TotalPrice:    0,
		DiscountPrice: 0,
		FinalPrice:    0,
		Status:        OrderStatus_Pending,
	}
}
