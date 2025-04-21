package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/category/dto/req"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type CategoryService interface {
	Create(req dtoreq.CreateCategoryReq) (*entities.Category, error)
	DeleteById(id uint) error
	GetCategoriesTree() ([]*entities.Category, error)
	GetCategories(page, pageSize int) ([]*entities.Category, int, error)
}

type CategoryServiceImpl struct {
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewCategoryServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
	}
}

func (svc CategoryServiceImpl) Create(dto dtoreq.CreateCategoryReq) (*entities.Category, error) {
	isNameDuplicated, duplicatedNameErr := svc.repo.CategoryRepo.Exist(map[string]any{
		"name": dto.Name,
	})
	if duplicatedNameErr != nil {
		return nil, types.NewServerError(
			"Error in checking category name exist",
			"CategoryServiceImpl.Create",
			duplicatedNameErr,
		)
	}
	if isNameDuplicated {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("category.errors.name_duplicate"),
		)
	}
	if dto.ParentID != nil {
		parentCategory, parentCategoryErr := svc.repo.CategoryRepo.GetByID(*dto.ParentID)
		if parentCategoryErr != nil {
			return nil, types.NewServerError(
				"Error in fetching parent category by id",
				"CategoryServiceImpl.Create",
				parentCategoryErr,
			)
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
	if err := svc.repo.CategoryRepo.Create(category); err != nil {
		return nil, types.NewServerError(
			"Create Category Throw Error",
			"CategoryServiceImpl.Create",
			err,
		)
	}
	return category, nil
}

func (svc CategoryServiceImpl) DeleteById(id uint) error {
	category, categoryErr := svc.repo.CategoryRepo.GetByID(id)
	if categoryErr != nil {
		return types.NewServerError(
			"Error in fetching category by id",
			"CategoryServiceImpl.DeleteById",
			categoryErr,
		)
	}
	if category == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("category.errors.not_found"),
		)
	}
	if err := svc.repo.CategoryRepo.Delete(category); err != nil {
		return types.NewServerError(
			"Delete Category By ID Throw Error",
			"CategoryServiceImpl.DeleteById",
			err,
		)
	}
	return nil
}

func (svc CategoryServiceImpl) getSubCategories(category *entities.Category) ([]*entities.Category, error) {
	subCategories, subCategoriesErr := svc.repo.CategoryRepo.GetChildren(category.ID)
	if subCategoriesErr != nil {
		return nil, types.NewServerError(
			"Error in fetching sub categories",
			"CategoryServiceImpl.getSubCategories",
			subCategoriesErr,
		)
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
	rootCategories, rootCategoriesErr := svc.repo.CategoryRepo.GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"parent_id": nil,
		},
	})
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

func (svc CategoryServiceImpl) GetCategories(page, pageSize int) ([]*entities.Category, int, error) {
	categories, count, categoriesErr := svc.repo.CategoryRepo.GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
	})
	if categoriesErr != nil {
		return nil, 0, types.NewServerError("Get Categories List Throw Error",
			"CategoryServiceImpl.FetchPage",
			categoriesErr,
		)
	}
	return categories, count, nil
}
