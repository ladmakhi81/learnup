package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/category/dto/req"
	categoryError "github.com/ladmakhi81/learnup/internals/category/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/types"
)

type CategoryService interface {
	Create(req dtoreq.CreateCategoryReq) (*entities.Category, error)
	DeleteById(id uint) error
	GetCategoriesTree() ([]*entities.Category, error)
	GetCategories(page, pageSize int) ([]*entities.Category, int, error)
}

type categoryService struct {
	unitOfWork db.UnitOfWork
}

func NewCategorySvc(unitOfWork db.UnitOfWork) CategoryService {
	return &categoryService{unitOfWork: unitOfWork}
}

func (svc categoryService) Create(dto dtoreq.CreateCategoryReq) (*entities.Category, error) {
	const operationName = "categoryService.Create"
	isNameDuplicated, err := svc.unitOfWork.CategoryRepo().Exist(map[string]any{"name": dto.Name})
	if err != nil {
		return nil, types.NewServerError("Error in checking category name exist", operationName, err)
	}
	if isNameDuplicated {
		return nil, categoryError.Category_DuplicateName
	}
	if dto.ParentID != nil {
		parentCategory, err := svc.unitOfWork.CategoryRepo().GetByID(*dto.ParentID, nil)
		if err != nil {
			return nil, types.NewServerError("Error in fetching parent category by id", operationName, err)
		}
		if parentCategory == nil {
			return nil, categoryError.Category_ParentNotFound
		}
	}
	category := &entities.Category{
		Name:             dto.Name,
		ParentCategoryID: dto.ParentID,
	}
	if err := svc.unitOfWork.CategoryRepo().Create(category); err != nil {
		return nil, types.NewServerError("Create Category Throw Error", operationName, err)
	}
	return category, nil
}

func (svc categoryService) DeleteById(id uint) error {
	const operationName = "categoryService.DeleteById"
	category, err := svc.unitOfWork.CategoryRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError("Error in fetching category by id", operationName, err)
	}
	if category == nil {
		return categoryError.Category_NotFound
	}
	if err := svc.unitOfWork.CategoryRepo().Delete(category); err != nil {
		return types.NewServerError("Delete Category By ID Throw Error", operationName, err)
	}
	return nil
}

func (svc categoryService) getSubCategories(category *entities.Category) ([]*entities.Category, error) {
	const operationName = "categoryService.getSubCategories"
	subCategories, err := svc.unitOfWork.CategoryRepo().GetChildren(category.ID)
	if err != nil {
		return nil, types.NewServerError("Error in fetching sub categories", operationName, err)
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

func (svc categoryService) GetCategoriesTree() ([]*entities.Category, error) {
	const operationName = "categoryService.GetCategoriesTree"
	rootCategories, err := svc.unitOfWork.CategoryRepo().GetAll(repositories.GetAllOptions{Conditions: map[string]any{"parent_category_id": nil}})
	if err != nil {
		return nil, types.NewServerError("Fetch Categories As Tree Throw Error", operationName, err)
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

func (svc categoryService) GetCategories(page, pageSize int) ([]*entities.Category, int, error) {
	const operationName = "categoryService.GetCategories"
	categories, count, err := svc.unitOfWork.CategoryRepo().GetPaginated(repositories.GetPaginatedOptions{Offset: &page, Limit: &pageSize})
	if err != nil {
		return nil, 0, types.NewServerError("Get Categories List Throw Error", operationName, err)
	}
	return categories, count, nil
}
