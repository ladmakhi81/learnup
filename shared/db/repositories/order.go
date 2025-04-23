package repositories

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
)

type OrderRepo interface {
	Repository[entities.Order]
}

type OrderRepoImpl struct {
	RepositoryImpl[entities.Order]
}

func NewOrderRepo(db *gorm.DB) *OrderRepoImpl {
	return &OrderRepoImpl{
		RepositoryImpl[entities.Order]{
			db: db,
		},
	}
}
