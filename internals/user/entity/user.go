package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model

	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Phone     string `gorm:"unique,column:phone_number"`
	Password  string `gorm:"not null,column:password"`
}

// Change default name of users table
func (User) TableName() string {
	return "_users"
}
