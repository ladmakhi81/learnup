package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/category/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/category/dto/res"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	categorySvc    categoryService.CategoryService
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
}

func NewHandler(
	categorySvc categoryService.CategoryService,
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
) *Handler {
	return &Handler{
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
//	@Success	201		{object}	types.ApiResponse{data=dtores.CreateCategoryResDto}
//	@Failure	400		{object}	types.ApiError
//	@Failure	404		{object}	types.ApiError
//	@Failure	409		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Router		/categories/ [post]
//
// @Security BearerAuth
func (h Handler) CreateCategory(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := new(dtoreq.CreateCategoryReq)
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	category, err := h.categorySvc.Create(*dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, dtores.NewCreateCategoryResDto(category)), nil
}

// GetCategoriesTree godoc
//
//	@Summary	Get categories as a tree structure
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	types.ApiResponse{data=[]dtores.GetCategoriesTreeItemDto}
//	@Failure	500	{object}	types.ApiResponse
//	@Router		/categories/tree [get]
//
// @Security BearerAuth
func (h Handler) GetCategoriesTree(_ *gin.Context) (*types.ApiResponse, error) {
	categories, err := h.categorySvc.GetCategoriesTree()
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, dtores.MapGetCategoriesTreeItemsDto(categories)), nil
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
//	@Router		/categories/page [get]
//
// @Security BearerAuth
func (h Handler) GetCategories(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	categories, count, err := h.categorySvc.GetCategories(page, pageSize)
	if err != nil {
		return nil, err
	}
	pageableCategoryRes := types.NewPaginationRes(
		dtores.MapCategoryPageableItemsDto(categories),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
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
//	@Router		/categories/{categoryId} [delete]
//
// @Security BearerAuth
func (h Handler) DeleteCategory(ctx *gin.Context) (*types.ApiResponse, error) {
	categoryID, err := utils.ToUint(ctx.Param("categoryId"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("category.errors.invalid_category_id"),
		)
	}
	if err := h.categorySvc.DeleteById(categoryID); err != nil {
		return nil, err
	}
	return types.NewApiResponse(
		http.StatusOK,
		map[string]any{},
	), nil
}
