package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type CategoryRepo interface {
	Create(category *entities.Category) error
	Delete(category *entities.Category) error
	FetchById(categoryID uint) (*entities.Category, error)
	FetchByName(name string) (*entities.Category, error)
	FetchPage(page, pageSize int) ([]*entities.Category, error)
	FetchCount() (int, error)
	FetchRoot() ([]*entities.Category, error)
	FetchChildren(categoryID uint) ([]*entities.Category, error)
}

type CategoryRepoImpl struct {
	db *db.Database
}

func NewCategoryRepoImpl(db *db.Database) *CategoryRepoImpl {
	return &CategoryRepoImpl{
		db: db,
	}
}

func (repo CategoryRepoImpl) Create(category *entities.Category) error {
	tx := repo.db.Core.Create(category)
	return tx.Error
}

func (repo CategoryRepoImpl) Delete(category *entities.Category) error {
	tx := repo.db.Core.Delete(category)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo CategoryRepoImpl) FetchById(categoryID uint) (*entities.Category, error) {
	category := &entities.Category{}
	tx := repo.db.Core.Where("id = ?", categoryID).First(category)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return category, nil
}

func (repo CategoryRepoImpl) FetchByName(name string) (*entities.Category, error) {
	category := &entities.Category{}
	tx := repo.db.Core.Where("name = ?", name).First(category)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return category, nil
}

func (repo CategoryRepoImpl) FetchRoot() ([]*entities.Category, error) {
	var categories []*entities.Category
	tx := repo.db.Core.Where("parent_category_id IS NULL").Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}

func (repo CategoryRepoImpl) FetchChildren(categoryID uint) ([]*entities.Category, error) {
	var categories []*entities.Category
	tx := repo.db.Core.
		Preload("Children").
		Where("parent_category_id = ?", categoryID).
		Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}

func (repo CategoryRepoImpl) FetchPage(page, pageSize int) ([]*entities.Category, error) {
	var categories []*entities.Category
	tx := repo.db.Core.
		Order("created_at desc").
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}

func (repo CategoryRepoImpl) FetchCount() (int, error) {
	count := int64(0)
	tx := repo.db.Core.Model(&entities.Category{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
