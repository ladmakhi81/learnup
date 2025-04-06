package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/category/dto/req"
	"github.com/ladmakhi81/learnup/internals/category/entity"
)

type CategoryService interface {
	Create(req dtoreq.CreateCategoryReq) (*entity.Category, error)
	Delete(id uint) error
	FindByID(id uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	IsCategoryNameExist(name string) (bool, error)
}

type CategoryServiceImpl struct{}

func NewCategoryServiceImpl() *CategoryServiceImpl {
	return &CategoryServiceImpl{}
}

func (service *CategoryServiceImpl) Create(req dtoreq.CreateCategoryReq) (*entity.Category, error) {
	return nil, nil
}

func (service *CategoryServiceImpl) Delete(id uint) error {
	return nil
}

func (service *CategoryServiceImpl) FindByID(id uint) (*entity.Category, error) {
	return nil, nil
}

func (service *CategoryServiceImpl) FindByName(name string) (*entity.Category, error) {
	return nil, nil
}

func (service *CategoryServiceImpl) IsCategoryNameExist(name string) (bool, error) {
	return false, nil
}
