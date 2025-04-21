package repositories

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type CategoryRepo interface {
	Repository[entities.Category]
	GetChildren(categoryID uint) ([]*entities.Category, error)
}

type CategoryRepoImpl struct {
	RepositoryImpl[entities.Category]
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepoImpl {
	return &CategoryRepoImpl{
		RepositoryImpl[entities.Category]{
			db: db,
		},
	}
}

func (repo CategoryRepoImpl) GetChildren(categoryID uint) ([]*entities.Category, error) {
	var categories []*entities.Category
	tx := repo.db.
		Preload("Children").
		Where("parent_category_id = ?", categoryID).
		Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}
