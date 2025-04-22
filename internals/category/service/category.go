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
	const operationName = "CategoryServiceImpl.Create"
	isNameDuplicated, err := svc.unitOfWork.CategoryRepo().Exist(map[string]any{
		"name": dto.Name,
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in checking category name exist",
			operationName,
			err,
		)
	}
	if isNameDuplicated {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("category.errors.name_duplicate"),
		)
	}
	if dto.ParentID != nil {
		parentCategory, err := svc.unitOfWork.CategoryRepo().GetByID(*dto.ParentID, nil)
		if err != nil {
			return nil, types.NewServerError(
				"Error in fetching parent category by id",
				operationName,
				err,
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
	if err := svc.unitOfWork.CategoryRepo().Create(category); err != nil {
		return nil, types.NewServerError(
			"Create Category Throw Error",
			operationName,
			err,
		)
	}
	return category, nil
}

func (svc CategoryServiceImpl) DeleteById(id uint) error {
	const operationName = "CategoryServiceImpl.DeleteById"
	category, err := svc.unitOfWork.CategoryRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError(
			"Error in fetching category by id",
			operationName,
			err,
		)
	}
	if category == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("category.errors.not_found"),
		)
	}
	if err := svc.unitOfWork.CategoryRepo().Delete(category); err != nil {
		return types.NewServerError(
			"Delete Category By ID Throw Error",
			operationName,
			err,
		)
	}
	return nil
}

func (svc CategoryServiceImpl) getSubCategories(category *entities.Category) ([]*entities.Category, error) {
	const operationName = "CategoryServiceImpl.getSubCategories"
	subCategories, err := svc.unitOfWork.CategoryRepo().GetChildren(category.ID)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching sub categories",
			operationName,
			err,
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
	const operationName = "CategoryServiceImpl.GetCategoriesTree"
	rootCategories, err := svc.unitOfWork.CategoryRepo().GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"parent_category_id": nil,
		},
	})
	if err != nil {
		return nil, types.NewServerError("Fetch Categories As Tree Throw Error",
			operationName,
			err,
		)
	}
	treeCategories := make([]*entities.Category, len(rootCategories))
	for rootCategoryIndex, rootCategory := range rootCategories {
		subCategory, err := svc.getSubCategories(rootCategory)
		if err != nil {
			return nil, err
		}
		rootCategory.Children = subCategory
		treeCategories[rootCategoryIndex] = rootCategory
	}
	return treeCategories, nil
}

func (svc CategoryServiceImpl) GetCategories(page, pageSize int) ([]*entities.Category, int, error) {
	const operationName = "CategoryServiceImpl.GetCategories"
	categories, count, err := svc.unitOfWork.CategoryRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
	})
	if err != nil {
		return nil, 0, types.NewServerError(
			"Get Categories List Throw Error",
			operationName,
			err,
		)
	}
	return categories, count, nil
}
