package res

import (
	"github.com/ladmakhi81/learnup/internals/user/entity"
	"time"
)

type CreateBasicUserRes struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewCreateUserResponse(user *entity.User) *CreateBasicUserRes {
	return &CreateBasicUserRes{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
