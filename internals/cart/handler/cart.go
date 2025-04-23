package handler

import (
	"github.com/gin-gonic/gin"
	cartDtoReq "github.com/ladmakhi81/learnup/internals/cart/dto/req"
	cartDtoRes "github.com/ladmakhi81/learnup/internals/cart/dto/res"
	cartService "github.com/ladmakhi81/learnup/internals/cart/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
	cartSvc        cartService.CartService
	userSvc        userService.UserSvc
}

func NewHandler(
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
	cartSvc cartService.CartService,
	userSvc userService.UserSvc,
) *Handler {
	return &Handler{
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
		cartSvc:        cartSvc,
		userSvc:        userSvc,
	}
}

// AddCart godoc
//
//	@Summary	Add a new cart item
//	@Tags		carts
//	@Accept		json
//	@Produce	json
//	@Param		body	body		cartDtoReq.CreateCartReqDto	true	" "
//	@Success	201		{object}	types.ApiResponse{data=cartDtoRes.AddCartResDto}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	404		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/carts [post]
func (h Handler) AddCart(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &cartDtoReq.CreateCartReqDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	cart, err := h.cartSvc.Create(user, *dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, cartDtoRes.NewAddCartResDto(cart)), nil
}

// DeleteCartByID godoc
//
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
	cartID, err := utils.ToUint(ctx.Param("cart-id"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("cart.errors.invalid_id"),
		)
	}
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.cartSvc.DeleteByID(user.ID, cartID); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

// GetCartsByUserID godoc
//
//	@Summary	Get all cart items for the authenticated user
//	@Tags		carts
//	@Produce	json
//	@Success	200	{object}	types.ApiResponse{data=[]cartDtoRes.GetCartItemDto}
//	@Failure	401	{object}	types.ApiError
//	@Failure	500	{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/carts [get]
func (h Handler) GetCartsByUserID(ctx *gin.Context) (*types.ApiResponse, error) {
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	carts, cartsErr := h.cartSvc.FetchAllByUserID(user.ID)
	if cartsErr != nil {
		return nil, cartsErr
	}
	cartRes := cartDtoRes.MapGetCartItemDto(carts)
	return types.NewApiResponse(http.StatusOK, cartRes), nil
}
