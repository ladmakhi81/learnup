package repo

import (
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/category/entity"
)

type CategoryRepo interface {
	Create(category *entity.Category) error
	Delete(categoryID uint) error
	FindByID(categoryID uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
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
	return nil
}

func (repo CategoryRepoImpl) Delete(categoryID uint) error {
	return nil
}

func (repo CategoryRepoImpl) FindByID(categoryID uint) (*entity.Category, error) {
	return nil, nil
}

func (repo CategoryRepoImpl) FindByName(name string) (*entity.Category, error) {
	return nil, nil
}
