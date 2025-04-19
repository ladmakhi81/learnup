package handler

import (
	"github.com/gin-gonic/gin"
	cartDtoReq "github.com/ladmakhi81/learnup/internals/cart/dto/req"
	cartDtoRes "github.com/ladmakhi81/learnup/internals/cart/dto/res"
	cartService "github.com/ladmakhi81/learnup/internals/cart/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
	cartSvc        cartService.CartService
}

func NewHandler(
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
	cartSvc cartService.CartService,
) *Handler {
	return &Handler{
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
		cartSvc:        cartSvc,
	}
}

// AddCart godoc
//
//	@Summary	Add a new cart item
//	@Tags		carts
//	@Accept		json
//	@Produce	json
//	@Param		body	body		cartDtoReq.CreateCartReq	true	" "
//	@Success	201		{object}	types.ApiResponse{data=cartDtoRes.AddCartRes}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	404		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/carts [post]
func (h Handler) AddCart(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &cartDtoReq.CreateCartReq{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	authContext, _ := ctx.Get("AUTH")
	authClaim := authContext.(*types.TokenClaim)
	dto.UserID = authClaim.UserID
	cart, cartErr := h.cartSvc.Create(*dto)
	if cartErr != nil {
		return nil, cartErr
	}
	cartRes := cartDtoRes.AddCartRes{
		ID:       cart.ID,
		UserID:   cart.UserID,
		CourseID: cart.CourseID,
	}
	return types.NewApiResponse(http.StatusCreated, cartRes), nil
}

// DeleteCartByID godoc
//	@Summary	Delete a cart item by ID
//	@Tags		carts
//	@Param		cart-id	path		uint	true	"Cart ID"
//	@Success	200		{object}	types.ApiResponse
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	404		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/carts/{cart-id} [delete]
func (h Handler) DeleteCartByID(ctx *gin.Context) (*types.ApiResponse, error) {
	parsedCartID, parsedErr := utils.ToUint(ctx.Param("cart-id"))
	if parsedErr != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("cart.errors.invalid_id"),
		)
	}
	authContext, _ := ctx.Get("AUTH")
	authClaim := authContext.(*types.TokenClaim)
	userID := authClaim.UserID
	if err := h.cartSvc.DeleteByID(userID, parsedCartID); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

// GetCartsByUserID godoc
//	@Summary	Get all cart items for the authenticated user
//	@Tags		carts
//	@Produce	json
//	@Success	200	{object}	types.ApiResponse{data=[]cartDtoRes.GetCartItem}
//	@Failure	401	{object}	types.ApiError
//	@Failure	500	{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/carts [get]
func (h Handler) GetCartsByUserID(ctx *gin.Context) (*types.ApiResponse, error) {
	authContext, _ := ctx.Get("AUTH")
	authClaim := authContext.(*types.TokenClaim)
	userID := authClaim.UserID
	carts, cartsErr := h.cartSvc.FetchAllByUserID(userID)
	if cartsErr != nil {
		return nil, cartsErr
	}
	cartRes := cartDtoRes.MapGetCartItems(carts)
	return types.NewApiResponse(http.StatusOK, cartRes), nil
}
