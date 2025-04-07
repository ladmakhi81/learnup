package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/category/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/category/dto/res"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/pkg/validation"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
	"strconv"
)

type CategoryAdminHandler struct {
	categorySvc    categoryService.CategoryService
	translationSvc translations.Translator
	validationSvc  validation.Validation
}

func NewCategoryAdminHandler(
	categorySvc categoryService.CategoryService,
	translationSvc translations.Translator,
	validationSvc validation.Validation,
) *CategoryAdminHandler {
	return &CategoryAdminHandler{
		categorySvc:    categorySvc,
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
	}
}

// CreateCategory godoc
//
//	@Summary	Create a new category
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtoreq.CreateCategoryReq	true	" "
//	@Success	201		{object}	types.ApiResponse{data=dtores.CreateCategoryRes}
//	@Failure	400		{object}	types.ApiError
//	@Failure	404		{object}	types.ApiError
//	@Failure	409		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Router		/categories/admin/ [post]
//
// @Security BearerAuth
func (h CategoryAdminHandler) CreateCategory(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := new(dtoreq.CreateCategoryReq)
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	category, categoryErr := h.categorySvc.Create(*dto)
	if categoryErr != nil {
		return nil, categoryErr
	}
	categoryRes := dtores.NewCreateCategoryRes(
		category.ID,
		category.Name,
		category.CreatedAt,
	)
	return types.NewApiResponse(http.StatusCreated, categoryRes), nil
}

// GetCategoriesTree godoc
//
//	@Summary	Get categories as a tree structure
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	types.ApiResponse{data=dtores.GetCategoriesTreeRes}
//	@Failure	500	{object}	types.ApiResponse
//	@Router		/categories/admin/tree [get]
//
// @Security BearerAuth
func (h CategoryAdminHandler) GetCategoriesTree(ctx *gin.Context) (*types.ApiResponse, error) {
	categoriesTree, categoriesTreeErr := h.categorySvc.GetCategoriesTree()
	if categoriesTreeErr != nil {
		return nil, categoriesTreeErr
	}
	//categoriesTreeRes := dtores.NewGetCategoriesTreeRes(categoriesTree)
	return types.NewApiResponse(http.StatusOK, categoriesTree), nil
}

// GetCategories godoc
//
//	@Summary	Get a paginated list of categories
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int											false	"Page number"				default(0)
//	@Param		pageSize	query		int											false	"Number of items per page"	default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes}	" "
//	@Failure	500			{object}	types.ApiError
//	@Router		/categories/admin/page [get]
//
// @Security BearerAuth
func (h CategoryAdminHandler) GetCategories(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	categories, categoriesErr := h.categorySvc.GetCategories(page, pageSize)
	if categoriesErr != nil {
		return nil, categoriesErr
	}
	categoriesCount, categoriesCountErr := h.categorySvc.GetCategoriesCount()
	if categoriesCountErr != nil {
		return nil, categoriesCountErr
	}
	pageableCategoryItems := dtores.MapCategoriesToPageableItems(categories)
	pageableCategoryRes := types.NewPaginationRes(
		pageableCategoryItems,
		page,
		utils.CalculatePaginationTotalPage(categoriesCount),
		categoriesCount,
	)

	return types.NewApiResponse(http.StatusOK, pageableCategoryRes), nil
}

// DeleteCategory godoc
//
//	@Summary	Delete a category by ID
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		categoryId	path		int	true	" "
//	@Success	200			{object}	types.ApiResponse
//	@Failure	400			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/categories/admin/{categoryId} [delete]
//
// @Security BearerAuth
func (h CategoryAdminHandler) DeleteCategory(ctx *gin.Context) (*types.ApiResponse, error) {
	categoryIDParam := ctx.Param("categoryId")
	categoryId, categoryIdErr := strconv.Atoi(categoryIDParam)
	if categoryIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("category.errors.invalid_category_id"),
		)
	}
	if err := h.categorySvc.DeleteById(uint(categoryId)); err != nil {
		return nil, err
	}
	return types.NewApiResponse(
		http.StatusOK,
		map[string]any{},
	), nil
}
