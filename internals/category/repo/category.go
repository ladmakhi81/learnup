package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/category/entity"
	"gorm.io/gorm"
)

type CategoryRepo interface {
	Create(category *entity.Category) error
	Delete(categoryID uint) error
	FindByID(categoryID uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	GetCategories(page, pageSize int) ([]*entity.Category, error)
	GetCategoriesCount() (int, error)
	GetRootCategories() ([]*entity.Category, error)
	GetSubCategories(categoryID uint) ([]*entity.Category, error)
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

func (repo CategoryRepoImpl) Delete(categoryID uint) error {
	return nil
}

func (repo CategoryRepoImpl) FindByID(categoryID uint) (*entity.Category, error) {
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

func (repo CategoryRepoImpl) FindByName(name string) (*entity.Category, error) {
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

func (repo CategoryRepoImpl) GetRootCategories() ([]*entity.Category, error) {
	var categories []*entity.Category
	tx := repo.db.Core.Where("parent_category_id IS NULL").Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}

func (repo CategoryRepoImpl) GetSubCategories(categoryID uint) ([]*entity.Category, error) {
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

func (repo CategoryRepoImpl) GetCategories(page, pageSize int) ([]*entity.Category, error) {
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

func (repo CategoryRepoImpl) GetCategoriesCount() (int, error) {
	count := int64(0)
	tx := repo.db.Core.Model(&entity.Category{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
