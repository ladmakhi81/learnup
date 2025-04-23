package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"time"
)

type courseCartItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetCartItemDto struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Course    *courseCartItem `json:"course"`
}

func MapGetCartItemDto(cartItems []*entities.Cart) []*GetCartItemDto {
	res := make([]*GetCartItemDto, len(cartItems))
	for index, cart := range cartItems {
		res[index] = &GetCartItemDto{
			ID:        cart.ID,
			CreatedAt: cart.CreatedAt,
			UpdatedAt: cart.UpdatedAt,
			Course: &courseCartItem{
				ID:          cart.Course.ID,
				Name:        cart.Course.Name,
				Description: cart.Course.Description,
			},
		}
	}
	return res
}
