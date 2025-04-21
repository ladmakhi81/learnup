package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type OrderRepo interface {
	Create(order *entities.Order) error
	Update(order *entities.Order) error
	FetchPaginated(page, pageSize int) ([]*entities.Order, error)
	FetchCount() (int, error)
	FetchDetailById(id uint) (*entities.Order, error)
}

type OrderRepoImpl struct {
	dbClient *db.Database
}

func NewOrderRepo(
	dbClient *db.Database,
) *OrderRepoImpl {
	return &OrderRepoImpl{
		dbClient: dbClient,
	}
}

func (repo OrderRepoImpl) Create(order *entities.Order) error {
	tx := repo.dbClient.Core.Create(order)
	return tx.Error
}

func (repo OrderRepoImpl) Update(order *entities.Order) error {
	tx := repo.dbClient.Core.Updates(order)
	return tx.Error
}

func (repo OrderRepoImpl) FetchPaginated(page, pageSize int) ([]*entities.Order, error) {
	var orders []*entities.Order
	tx := repo.dbClient.Core.
		Preload("User").
		Offset(page * pageSize).
		Limit(pageSize).
		Order("created_at desc").
		Find(&orders)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return orders, nil
}

func (repo OrderRepoImpl) FetchCount() (int, error) {
	var count int64
	tx := repo.dbClient.Core.Model(&entities.Order{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}

func (repo OrderRepoImpl) FetchDetailById(id uint) (*entities.Order, error) {
	var order *entities.Order
	tx := repo.dbClient.Core.
		Where("id = ?", id).
		Preload("User").
		Preload("Items").
		Preload("Items.Course").
		First(&order)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return order, nil
}
