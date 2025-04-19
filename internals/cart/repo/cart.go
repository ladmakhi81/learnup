package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type CartRepo interface {
	Create(cart *entities.Cart) error
	FetchByID(id uint) (*entities.Cart, error)
	FetchAllByUserID(userID uint) ([]*entities.Cart, error)
	FetchByUserAndCourse(userID uint, courseID uint) (*entities.Cart, error)
	DeleteByID(id uint) error
	DeleteAllByUserID(userID uint) error
	FetchByCartIDs(ids []uint) ([]*entities.Cart, error)
}

type CartRepoImpl struct {
	dbClient *db.Database
}

func NewCartRepo(
	dbClient *db.Database,
) *CartRepoImpl {
	return &CartRepoImpl{
		dbClient: dbClient,
	}
}

func (repo CartRepoImpl) Create(cart *entities.Cart) error {
	tx := repo.dbClient.Core.Create(cart)
	return tx.Error
}

func (repo CartRepoImpl) FetchByID(id uint) (*entities.Cart, error) {
	var cart *entities.Cart
	tx := repo.dbClient.Core.Where("id = ?", id).First(&cart)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return cart, nil
}

func (repo CartRepoImpl) DeleteByID(id uint) error {
	tx := repo.dbClient.Core.Where("id = ?", id).Delete(&entities.Cart{})
	return tx.Error
}

func (repo CartRepoImpl) DeleteAllByUserID(userID uint) error {
	tx := repo.dbClient.Core.Where("user_id = ?", userID).Delete(&entities.Cart{})
	return tx.Error
}

func (repo CartRepoImpl) FetchAllByUserID(userID uint) ([]*entities.Cart, error) {
	var carts []*entities.Cart
	tx := repo.dbClient.Core.
		Where("user_id = ?", userID).
		Preload("Course").
		Order("created_at desc").
		Find(&carts)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return carts, nil
}

func (repo CartRepoImpl) FetchByUserAndCourse(userID uint, courseID uint) (*entities.Cart, error) {
	var cart *entities.Cart
	tx := repo.dbClient.Core.Where("user_id = ? AND course_id = ?", userID, courseID).First(&cart)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return cart, nil
}

func (repo CartRepoImpl) FetchByCartIDs(ids []uint) ([]*entities.Cart, error) {
	var carts []*entities.Cart
	tx := repo.dbClient.Core.
		Where("id IN (?)", ids).
		Preload("Course").
		Find(&carts)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return carts, nil
}
