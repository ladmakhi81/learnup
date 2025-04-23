package entities

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	OrderID       uint           `gorm:"column:order_id;type:int;index;not null;"`
	Order         *Order         `gorm:"foreignkey:order_id"`
	UserID        uint           `gorm:"column:user_id;type:int;index;not null;"`
	User          *User          `gorm:"foreignkey:user_id;"`
	Gateway       PaymentGateway `gorm:"column:gateway;type:varchar(255);not null;"`
	MerchantID    string         `gorm:"column:merchant_id;type:text;not null"`
	Status        PaymentStatus  `gorm:"column:status;type:varchar(255);not null;"`
	Authority     string         `gorm:"column:authority;type:text;not null;"`
	Amount          float64        `gorm:"column:amount;type:decimal(10,2);not null;"`
	StatusChangedAt *time.Time     `gorm:"column:status_changed_at;type:timestamp;"`
	PayLink         string         `gorm:"column:pay_link;type:text;not null;"`
	RefID           string         `gorm:"column:ref_id;type:text;default null;"`
	TransactionID *uint          `gorm:"column:transaction_id;type:int;index;default:null"`
	Transaction   *Transaction   `gorm:"foreignkey:transaction_id"`
}

func (Payment) TableName() string {
	return "_payments"
}
