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
	GetCategoriesTree() ([]*entity.Category, error)
	GetCategories(page, pageSize int) ([]*entity.Category, error)
	GetCategoriesCount() (int, error)
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

func (svc CategoryServiceImpl) getSubCategories(category *entity.Category) ([]*entity.Category, error) {
	subCategories, subCategoriesErr := svc.repo.GetSubCategories(category.ID)
	if subCategoriesErr != nil {
		return nil, subCategoriesErr
	}
	for _, subCategory := range subCategories {
		nextSubCategory, nextSubCategoryErr := svc.getSubCategories(subCategory)
		if nextSubCategoryErr != nil {
			return nil, nextSubCategoryErr
		}
		subCategory.Children = nextSubCategory
	}
	return subCategories, nil
}

func (svc CategoryServiceImpl) GetCategoriesTree() ([]*entity.Category, error) {
	rootCategories, rootCategoriesErr := svc.repo.GetRootCategories()
	if rootCategoriesErr != nil {
		return nil, types.NewServerError("Fetch Categories As Tree Throw Error",
			"CategoryServiceImpl.GetCategoriesTree",
			rootCategoriesErr,
		)
	}
	treeCategories := make([]*entity.Category, len(rootCategories))
	for rootCategoryIndex, rootCategory := range rootCategories {
		subCategory, subCategoryErr := svc.getSubCategories(rootCategory)
		if subCategoryErr != nil {
			return nil, subCategoryErr
		}
		rootCategory.Children = subCategory
		treeCategories[rootCategoryIndex] = rootCategory
	}
	return treeCategories, nil
}

func (svc CategoryServiceImpl) GetCategories(page, pageSize int) ([]*entity.Category, error) {
	categories, categoriesErr := svc.repo.GetCategories(page, pageSize)
	if categoriesErr != nil {
		return nil, types.NewServerError("Get Categories List Throw Error",
			"CategoryServiceImpl.GetCategories",
			categoriesErr,
		)
	}
	return categories, nil
}

func (svc CategoryServiceImpl) GetCategoriesCount() (int, error) {
	count, countErr := svc.repo.GetCategoriesCount()
	if countErr != nil {
		return 0, types.NewServerError("Get Count Of Categories Throw Error",
			"CategoryServiceImpl.GetCategoriesCount",
			countErr,
		)
	}
	return count, nil
}
