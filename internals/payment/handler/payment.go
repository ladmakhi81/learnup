package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/db/entities"
	paymentDtoReq "github.com/ladmakhi81/learnup/internals/payment/dto/req"
	paymentDtoRes "github.com/ladmakhi81/learnup/internals/payment/dto/res"
	paymentService "github.com/ladmakhi81/learnup/internals/payment/service"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	paymentSvc paymentService.PaymentService
}

func NewHandler(
	paymentSvc paymentService.PaymentService,
) *Handler {
	return &Handler{
		paymentSvc: paymentSvc,
	}
}

func (h Handler) VerifyZarinpal(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := paymentDtoReq.VerifyPaymentReq{
		Authority: ctx.Query("Authority"),
		Gateway:   entities.PaymentGateway_Zarinpal,
	}
	if err := h.paymentSvc.Verify(dto); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

func (h Handler) VerifyZibal(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := paymentDtoReq.VerifyPaymentReq{
		Authority: ctx.Query("trackId"),
		Gateway:   entities.PaymentGateway_Zibal,
	}
	if err := h.paymentSvc.Verify(dto); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

func (h Handler) VerifyStripe(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := paymentDtoReq.VerifyPaymentReq{
		Authority: ctx.Query("session_id"),
		Gateway:   entities.PaymentGateway_Stripe,
	}
	if err := h.paymentSvc.Verify(dto); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

// GetPayments godoc
//
//	@Summary	Get paginated payments
//	@Tags		payments
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=paymentDtoRes.GetPageablePaymentItem}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/payments/page [get]
//	@Security	BearerAuth
func (h Handler) GetPayments(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	payments, paymentsErr := h.paymentSvc.FetchPageable(page, pageSize)
	if paymentsErr != nil {
		return nil, paymentsErr
	}
	paymentCount, paymentCountErr := h.paymentSvc.FetchCount()
	if paymentCountErr != nil {
		return nil, paymentCountErr
	}
	res := types.NewPaginationRes(
		paymentDtoRes.MapGetPageablePaymentItems(payments),
		page,
		utils.CalculatePaginationTotalPage(paymentCount, pageSize),
		paymentCount,
	)
	return types.NewApiResponse(http.StatusOK, res), nil
}
