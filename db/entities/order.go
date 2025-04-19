package entities

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID          uint         `gorm:"column:user_id;type:int;not null;index"`
	User            *User        `gorm:"foreignkey:user_id"`
	FinalAmount     float64      `gorm:"column:final_amount;type:decimal(10,2);not null"`
	DiscountAmount  float64      `gorm:"column:discount_amount;type:decimal(10,2);default:0"`
	Amount          float64      `gorm:"column:amount;type:decimal(10,2);not null"`
	Status          OrderStatus  `gorm:"column:status;type:varchar(255);default:'pending';"`
	StatusChangedAt *time.Time   `gorm:"column:status_changed_at;type:timestamp;default:null"`
	Items           []*OrderItem `gorm:"foreignkey:order_id"`
}

func (Order) TableName() string {
	return "_orders"
}
