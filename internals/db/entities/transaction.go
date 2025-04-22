package entities

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model

	Amount   float64         `gorm:"column:amount;type:decimal(10,2);not null;"`
	User     string          `gorm:"column:user;type:varchar(255);not null;"`
	Phone    string          `gorm:"column:phone;type:varchar(255);not null;"`
	Type     TransactionType `gorm:"column:type;type:varchar(255);not null;"`
	Tag      TransactionTag  `gorm:"column:tag;type:varchar(255);not null;"`
	Currency string          `gorm:"column:currency;type:varchar(255);not null;"`
}

func (Transaction) TableName() string {
	return "_transactions"
}
