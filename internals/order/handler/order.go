package handler

import (
	"github.com/gin-gonic/gin"
	orderDtoReq "github.com/ladmakhi81/learnup/internals/order/dto/req"
	orderDtoRes "github.com/ladmakhi81/learnup/internals/order/dto/res"
	orderService "github.com/ladmakhi81/learnup/internals/order/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	orderSvc       orderService.OrderService
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
	userSvc        userService.UserSvc
}

func NewHandler(
	orderSvc orderService.OrderService,
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
	userSvc userService.UserSvc,
) *Handler {
	return &Handler{
		orderSvc:       orderSvc,
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
		userSvc:        userSvc,
	}
}

// CreateOrder godoc
//
//	@Summary	Create a new order
//	@Tags		orders
//	@Accept		json
//	@Produce	json
//	@Param		order	body		orderDtoReq.CreateOrderReqDto	true	" "
//	@Success	201		{object}	types.ApiResponse{data=orderDtoRes.CreateOrderResDto}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/orders [post]
func (h Handler) CreateOrder(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &orderDtoReq.CreateOrderReqDto{}
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
	payLink, err := h.orderSvc.Create(user, *dto)
	if err != nil {
		return nil, err
	}
	res := orderDtoRes.CreateOrderResDto{PayLink: payLink}
	return types.NewApiResponse(http.StatusCreated, res), nil
}

// GetOrdersPage godoc
//
//	@Summary	Retrieve paginated orders
//	@Tags		orders
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"		default(0)
//	@Param		pageSize	query		int	false	"Items per page"	default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]orderDtoRes.PaginatedOrderItemDto}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/orders [get]
func (h Handler) GetOrdersPage(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	orders, count, err := h.orderSvc.FetchPaginated(page, pageSize)
	if err != nil {
		return nil, err
	}
	res := types.NewPaginationRes(
		orderDtoRes.MapPaginatedOrderItemsDto(orders),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, res), nil
}

// GetOrderByID godoc
//
//	@Summary	Retrieve order details by ID
//	@Tags		orders
//	@Accept		json
//	@Produce	json
//	@Param		order-id	path		uint	true	" "
//	@Success	200			{object}	orderDtoRes.GetOrderDetailItemDto
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/orders/{order-id} [get]
func (h Handler) GetOrderByID(ctx *gin.Context) (*types.ApiResponse, error) {
	orderID, err := utils.ToUint(ctx.Param("order-id"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("order.errors.invalid_id"),
		)
	}
	order, err := h.orderSvc.FetchDetailById(orderID)
	if err != nil {
		return nil, err
	}
	res := orderDtoRes.NewGetOrderDetailItemDto(order)
	return types.NewApiResponse(http.StatusOK, res), nil
}
