package repo

import (
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
)

type OrderItemRepo interface {
	CreateBatch(orderItems []*entities.OrderItem) error
}

type OrderItemRepoImpl struct {
	dbClient *db.Database
}

func NewOrderItemRepo(dbClient *db.Database) *OrderItemRepoImpl {
	return &OrderItemRepoImpl{
		dbClient: dbClient,
	}
}

func (repo OrderItemRepoImpl) CreateBatch(orderItems []*entities.OrderItem) error {
	tx := repo.dbClient.Core.Create(&orderItems)
	return tx.Error
}
