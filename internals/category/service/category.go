package service

import (
	"fmt"
	dtoreq "github.com/ladmakhi81/learnup/internals/category/dto/req"
	"github.com/ladmakhi81/learnup/internals/category/entity"
	"github.com/ladmakhi81/learnup/internals/category/repo"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/types"
)

type CategoryService interface {
	Create(req dtoreq.CreateCategoryReq) (*entity.Category, error)
	Delete(id uint) error
	FindByID(id uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	IsCategoryNameExist(name string) (bool, error)
}

type CategoryServiceImpl struct {
	repo           repo.CategoryRepo
	translationSvc translations.Translator
}

func NewCategoryServiceImpl(
	repo repo.CategoryRepo,
	translationSvc translations.Translator,
) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
	}
}

func (svc CategoryServiceImpl) Create(dto dtoreq.CreateCategoryReq) (*entity.Category, error) {
	isDuplicateName, duplicateNameErr := svc.IsCategoryNameExist(dto.Name)
	if duplicateNameErr != nil {
		return nil, duplicateNameErr
	}
	if isDuplicateName {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("category.errors.name_duplicate"),
		)
	}
	fmt.Println("parent id", dto.ParentID, *dto.ParentID)
	if dto.ParentID != nil {
		parentCategory, parentCategoryErr := svc.FindByID(*dto.ParentID)
		if parentCategoryErr != nil {
			return nil, parentCategoryErr
		}
		if parentCategory == nil {
			return nil, types.NewNotFoundError(
				svc.translationSvc.Translate(
					"category.errors.parent_category_id_not_found",
				),
			)
		}
	}
	category := &entity.Category{
		Name:             dto.Name,
		ParentCategoryID: dto.ParentID,
	}
	if err := svc.repo.Create(category); err != nil {
		return nil, types.NewServerError(
			"Create Category Throw Error",
			"CategoryServiceImpl.Create",
			err,
		)
	}
	return category, nil
}

func (svc CategoryServiceImpl) Delete(id uint) error {
	return nil
}

func (svc CategoryServiceImpl) FindByID(id uint) (*entity.Category, error) {
	category, categoryErr := svc.repo.FindByID(id)
	if categoryErr != nil {
		return nil, types.NewServerError("Find Category By Id Throw Error",
			"CategoryServiceImpl.FindById",
			categoryErr,
		)
	}
	return category, nil
}

func (svc CategoryServiceImpl) FindByName(name string) (*entity.Category, error) {
	category, categoryErr := svc.repo.FindByName(name)
	if categoryErr != nil {
		return nil, types.NewServerError("Find Category By Name Throw Error",
			"CategoryServiceImpl.FindByName",
			categoryErr,
		)
	}
	return category, nil
}

func (svc CategoryServiceImpl) IsCategoryNameExist(name string) (bool, error) {
	category, categoryErr := svc.FindByName(name)
	if categoryErr != nil {
		return false, categoryErr
	}
	if category == nil {
		return false, nil
	}
	return true, nil
}
