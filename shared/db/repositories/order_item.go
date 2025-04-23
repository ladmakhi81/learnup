package repositories

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
)

type OrderItemRepo interface {
	Repository[entities.OrderItem]
}

type OrderItemRepoImpl struct {
	RepositoryImpl[entities.OrderItem]
}

func NewOrderItemRepo(db *gorm.DB) *OrderItemRepoImpl {
	return &OrderItemRepoImpl{
		RepositoryImpl[entities.OrderItem]{
			db: db,
		},
	}
}
