package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
	"time"
)

type userOrderItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type PaginatedOrderItemDto struct {
	ID              uint                  `json:"id"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt        `json:"deletedAt"`
	User            userOrderItem         `json:"user"`
	FinalPrice      float64               `json:"finalPrice"`
	DiscountPrice   float64               `json:"discountPrice"`
	TotalPrice      float64               `json:"totalPrice"`
	Status          entities2.OrderStatus `json:"status"`
	StatusChangedAt *time.Time            `json:"statusChangedAt"`
}

func MapPaginatedOrderItemsDto(orders []*entities2.Order) []*PaginatedOrderItemDto {
	res := make([]*PaginatedOrderItemDto, len(orders))
	for i, order := range orders {
		res[i] = &PaginatedOrderItemDto{
			ID:              order.ID,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
			TotalPrice:      order.TotalPrice,
			Status:          order.Status,
			FinalPrice:      order.FinalPrice,
			DiscountPrice:   order.DiscountPrice,
			DeletedAt:       order.DeletedAt,
			StatusChangedAt: order.StatusChangedAt,
			User: userOrderItem{
				ID:       order.User.ID,
				Phone:    order.User.Phone,
				FullName: order.User.FullName(),
			},
		}
	}
	return res
}
