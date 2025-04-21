package handler

import (
	"github.com/gin-gonic/gin"
	orderDtoReq "github.com/ladmakhi81/learnup/internals/order/dto/req"
	orderDtoRes "github.com/ladmakhi81/learnup/internals/order/dto/res"
	orderService "github.com/ladmakhi81/learnup/internals/order/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	orderSvc       orderService.OrderService
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
}

func NewHandler(
	orderSvc orderService.OrderService,
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
) *Handler {
	return &Handler{
		orderSvc:       orderSvc,
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
	}
}

// CreateOrder godoc
//
//	@Summary	Create a new order
//	@Tags		orders
//	@Accept		json
//	@Produce	json
//	@Param		order	body		orderDtoReq.CreateOrderReq	true	" "
//	@Success	201		{object}	types.ApiResponse{data=orderDtoRes.CreateOrderRes}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/orders [post]
func (h Handler) CreateOrder(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &orderDtoReq.CreateOrderReq{}
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
	userID := authClaim.UserID
	dto.UserID = userID
	payLink, err := h.orderSvc.Create(*dto)
	if err != nil {
		return nil, err
	}
	res := orderDtoRes.CreateOrderRes{PayLink: payLink}
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
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]orderDtoRes.PaginatedOrderItem}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/orders [get]
func (h Handler) GetOrdersPage(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	orders, count, ordersErr := h.orderSvc.FetchPaginated(page, pageSize)
	if ordersErr != nil {
		return nil, ordersErr
	}
	res := types.NewPaginationRes(
		orderDtoRes.MapPaginatedOrderItems(orders),
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
//	@Success	200			{object}	orderDtoRes.GetOrderDetailRes
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/orders/{order-id} [get]
func (h Handler) GetOrderByID(ctx *gin.Context) (*types.ApiResponse, error) {
	parsedOrderID, parsedErr := utils.ToUint(ctx.Param("order-id"))
	if parsedErr != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("order.errors.invalid_id"),
		)
	}
	order, orderErr := h.orderSvc.FetchDetailById(parsedOrderID)
	if orderErr != nil {
		return nil, orderErr
	}
	res := orderDtoRes.NewGetOrderDetailRes(order)
	return types.NewApiResponse(http.StatusOK, res), nil
}
