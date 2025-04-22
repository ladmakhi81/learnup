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
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewCategoryServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc CategoryServiceImpl) Create(dto dtoreq.CreateCategoryReq) (*entities.Category, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	isNameDuplicated, duplicatedNameErr := tx.CategoryRepo().Exist(map[string]any{
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
		parentCategory, parentCategoryErr := tx.CategoryRepo().GetByID(*dto.ParentID, nil)
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
	if err := tx.CategoryRepo().Create(category); err != nil {
		return nil, types.NewServerError(
			"Create Category Throw Error",
			"CategoryServiceImpl.Create",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return category, nil
}

func (svc CategoryServiceImpl) DeleteById(id uint) error {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return txErr
	}
	category, categoryErr := tx.CategoryRepo().GetByID(id, nil)
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
	if err := tx.CategoryRepo().Delete(category); err != nil {
		return types.NewServerError(
			"Delete Category By ID Throw Error",
			"CategoryServiceImpl.DeleteById",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (svc CategoryServiceImpl) getSubCategories(category *entities.Category) ([]*entities.Category, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	subCategories, subCategoriesErr := tx.CategoryRepo().GetChildren(category.ID)
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
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return subCategories, nil
}

func (svc CategoryServiceImpl) GetCategoriesTree() ([]*entities.Category, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	rootCategories, rootCategoriesErr := tx.CategoryRepo().GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"parent_category_id": nil,
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
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return treeCategories, nil
}

func (svc CategoryServiceImpl) GetCategories(page, pageSize int) ([]*entities.Category, int, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, 0, txErr
	}
	categories, count, categoriesErr := tx.CategoryRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
	})
	if categoriesErr != nil {
		return nil, 0, types.NewServerError("Get Categories List Throw Error",
			"CategoryServiceImpl.FetchPage",
			categoriesErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return categories, count, nil
}
