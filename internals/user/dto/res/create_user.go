package dtores

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type CreateBasicUserResDto struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewCreateBasicUserResDto(user *entities.User) CreateBasicUserResDto {
	return CreateBasicUserResDto{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
