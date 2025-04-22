package entities

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Phone     string `gorm:"index;column:phone_number"`
	Password  string `gorm:"not null;column:password"`
}

// Change default name of users table
func (User) TableName() string {
	return "_users"
}

func (user User) FullName() string {
	return user.FirstName + " " + user.LastName
}

func (user User) IsPasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
