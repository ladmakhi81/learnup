package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/category/entity"
	"gorm.io/gorm"
)

type CategoryRepo interface {
	Create(category *entity.Category) error
	Delete(category *entity.Category) error
	FetchById(categoryID uint) (*entity.Category, error)
	FetchByName(name string) (*entity.Category, error)
	FetchPage(page, pageSize int) ([]*entity.Category, error)
	FetchCount() (int, error)
	FetchRoot() ([]*entity.Category, error)
	FetchChildren(categoryID uint) ([]*entity.Category, error)
}

type CategoryRepoImpl struct {
	db *db.Database
}

func NewCategoryRepoImpl(db *db.Database) *CategoryRepoImpl {
	return &CategoryRepoImpl{
		db: db,
	}
}

func (repo CategoryRepoImpl) Create(category *entity.Category) error {
	tx := repo.db.Core.Create(category)
	return tx.Error
}

func (repo CategoryRepoImpl) Delete(category *entity.Category) error {
	tx := repo.db.Core.Delete(category)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo CategoryRepoImpl) FetchById(categoryID uint) (*entity.Category, error) {
	category := &entity.Category{}
	tx := repo.db.Core.Where("id = ?", categoryID).First(category)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return category, nil
}

func (repo CategoryRepoImpl) FetchByName(name string) (*entity.Category, error) {
	category := &entity.Category{}
	tx := repo.db.Core.Where("name = ?", name).First(category)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return category, nil
}

func (repo CategoryRepoImpl) FetchRoot() ([]*entity.Category, error) {
	var categories []*entity.Category
	tx := repo.db.Core.Where("parent_category_id IS NULL").Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}

func (repo CategoryRepoImpl) FetchChildren(categoryID uint) ([]*entity.Category, error) {
	var categories []*entity.Category
	tx := repo.db.Core.
		Preload("Children").
		Where("parent_category_id = ?", categoryID).
		Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}

func (repo CategoryRepoImpl) FetchPage(page, pageSize int) ([]*entity.Category, error) {
	var categories []*entity.Category
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
	tx := repo.db.Core.Model(&entity.Category{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
