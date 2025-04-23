package dtores

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type AddCartResDto struct {
	UserID    uint      `json:"userId"`
	CourseID  uint      `json:"courseId"`
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAddCartResDto(cart *entities.Cart) AddCartResDto {
	return AddCartResDto{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CourseID:  cart.CourseID,
		CreatedAt: cart.CreatedAt,
	}
}
