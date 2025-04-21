package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"gorm.io/gorm"
	"time"
)

type userOrderDetailItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type orderCourseItem struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type orderItem struct {
	ID     uint            `json:"id"`
	Amount float64         `json:"amount"`
	Course orderCourseItem `json:"course"`
}

type GetOrderDetailRes struct {
	ID              uint                  `json:"id"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt        `json:"deletedAt"`
	User            userOrderDetailItem   `json:"user"`
	FinalPrice      float64               `json:"finalPrice"`
	DiscountPrice   float64               `json:"discountPrice"`
	TotalPrice      float64               `json:"totalPrice"`
	Status          entities2.OrderStatus `json:"status"`
	StatusChangedAt *time.Time            `json:"statusChangedAt"`
	Items           []orderItem           `json:"items"`
}

func NewGetOrderDetailRes(order *entities2.Order) *GetOrderDetailRes {
	items := make([]orderItem, len(order.Items))

	for i, item := range order.Items {
		items[i] = orderItem{
			ID:     item.ID,
			Amount: item.Amount,
			Course: orderCourseItem{
				ID:          item.Course.ID,
				Name:        item.Course.Name,
				Description: item.Course.Description,
				Price:       item.Course.Price,
			},
		}
	}

	return &GetOrderDetailRes{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,
		User: userOrderDetailItem{
			ID:       order.User.ID,
			FullName: order.User.FullName(),
			Phone:    order.User.Phone,
		},
		TotalPrice:      order.TotalPrice,
		DiscountPrice:   order.DiscountPrice,
		FinalPrice:      order.FinalPrice,
		Status:          order.Status,
		StatusChangedAt: order.StatusChangedAt,
		Items:           items,
	}
}
