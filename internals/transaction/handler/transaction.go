package handler

import (
	"github.com/gin-gonic/gin"
	transactionDtoRes "github.com/ladmakhi81/learnup/internals/transaction/dto/res"
	transactionService "github.com/ladmakhi81/learnup/internals/transaction/service"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	transactionSvc transactionService.TransactionService
}

func NewHandler(
	transactionSvc transactionService.TransactionService,
) *Handler {
	return &Handler{
		transactionSvc: transactionSvc,
	}
}

// GetTransactionsPage godoc
//
//	@Summary	Get paginated transactions
//	@Tags		transactions
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"		default(0)
//	@Param		pageSize	query		int	false	"Items per page"	default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=transactionDtoRes.GetTransactionPageableItem}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/transactions/page [get]
//	@Security	BearerAuth
func (h Handler) GetTransactionsPage(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	transactions, transactionErr := h.transactionSvc.FetchPageable(page, pageSize)
	if transactionErr != nil {
		return nil, transactionErr
	}
	count, countErr := h.transactionSvc.FetchCount()
	if countErr != nil {
		return nil, countErr
	}
	res := types.NewPaginationRes(
		transactionDtoRes.MapGetTransactionPageableItems(transactions),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, res), nil
}
