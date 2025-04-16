package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/category/dto/req"
	"github.com/ladmakhi81/learnup/internals/category/repo"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type CategoryService interface {
	Create(req dtoreq.CreateCategoryReq) (*entities.Category, error)
	DeleteById(id uint) error
	FindByID(id uint) (*entities.Category, error)
	FindByName(name string) (*entities.Category, error)
	IsCategoryNameExist(name string) (bool, error)
	GetCategoriesTree() ([]*entities.Category, error)
	GetCategories(page, pageSize int) ([]*entities.Category, error)
	GetCategoriesCount() (int, error)
}

type CategoryServiceImpl struct {
	repo           repo.CategoryRepo
	translationSvc contracts.Translator
}

func NewCategoryServiceImpl(
	repo repo.CategoryRepo,
	translationSvc contracts.Translator,
) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
	}
}

func (svc CategoryServiceImpl) Create(dto dtoreq.CreateCategoryReq) (*entities.Category, error) {
	isDuplicateName, duplicateNameErr := svc.IsCategoryNameExist(dto.Name)
	if duplicateNameErr != nil {
		return nil, duplicateNameErr
	}
	if isDuplicateName {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("category.errors.name_duplicate"),
		)
	}
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
	category := &entities.Category{
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

func (svc CategoryServiceImpl) DeleteById(id uint) error {
	category, categoryErr := svc.FindByID(id)
	if categoryErr != nil {
		return categoryErr
	}
	if category == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("category.errors.not_found"),
		)
	}
	if err := svc.repo.Delete(category); err != nil {
		return types.NewServerError(
			"Delete Category By ID Throw Error",
			"CategoryServiceImpl.DeleteById",
			err,
		)
	}
	return nil
}

func (svc CategoryServiceImpl) FindByID(id uint) (*entities.Category, error) {
	category, categoryErr := svc.repo.FetchById(id)
	if categoryErr != nil {
		return nil, types.NewServerError("Find Category By Id Throw Error",
			"CategoryServiceImpl.FetchById",
			categoryErr,
		)
	}
	return category, nil
}

func (svc CategoryServiceImpl) FindByName(name string) (*entities.Category, error) {
	category, categoryErr := svc.repo.FetchByName(name)
	if categoryErr != nil {
		return nil, types.NewServerError("Find Category By Name Throw Error",
			"CategoryServiceImpl.FetchByName",
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

func (svc CategoryServiceImpl) getSubCategories(category *entities.Category) ([]*entities.Category, error) {
	subCategories, subCategoriesErr := svc.repo.FetchChildren(category.ID)
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

func (svc CategoryServiceImpl) GetCategoriesTree() ([]*entities.Category, error) {
	rootCategories, rootCategoriesErr := svc.repo.FetchRoot()
	if rootCategoriesErr != nil {
		return nil, types.NewServerError("Fetch Categories As Tree Throw Error",
			"CategoryServiceImpl.GetCategoriesTree",
			rootCategoriesErr,
		)
	}
	treeCategories := make([]*entities.Category, len(rootCategories))
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

func (svc CategoryServiceImpl) GetCategories(page, pageSize int) ([]*entities.Category, error) {
	categories, categoriesErr := svc.repo.FetchPage(page, pageSize)
	if categoriesErr != nil {
		return nil, types.NewServerError("Get Categories List Throw Error",
			"CategoryServiceImpl.FetchPage",
			categoriesErr,
		)
	}
	return categories, nil
}

func (svc CategoryServiceImpl) GetCategoriesCount() (int, error) {
	count, countErr := svc.repo.FetchCount()
	if countErr != nil {
		return 0, types.NewServerError("Get Count Of Categories Throw Error",
			"CategoryServiceImpl.FetchCount",
			countErr,
		)
	}
	return count, nil
}
