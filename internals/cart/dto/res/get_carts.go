package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"time"
)

type courseCartItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetCartItem struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Course    *courseCartItem `json:"course"`
}

func MapGetCartItems(cartItems []*entities.Cart) []*GetCartItem {
	res := make([]*GetCartItem, len(cartItems))
	for index, cart := range cartItems {
		res[index] = &GetCartItem{
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
