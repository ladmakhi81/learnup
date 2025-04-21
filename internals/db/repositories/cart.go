package repositories

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"gorm.io/gorm"
)

type CartRepo interface {
	Repository[entities.Cart]
	GetByCartIDs(cartIDs []uint) ([]*entities.Cart, error)
}

type CartRepoImpl struct {
	RepositoryImpl[entities.Cart]
}

func NewCartRepo(db *gorm.DB) *CartRepoImpl {
	return &CartRepoImpl{
		RepositoryImpl[entities.Cart]{
			db: db,
		},
	}
}

func (repo CartRepoImpl) GetByCartIDs(cartIDs []uint) ([]*entities.Cart, error) {
	var carts []*entities.Cart
	tx := repo.db.
		Where("id IN (?)", cartIDs).
		Preload("Course").
		Find(&carts)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return carts, nil
}
