package handler

import (
	"github.com/gin-gonic/gin"
	dtores "github.com/ladmakhi81/learnup/internals/comment/dto/res"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/types"
	"github.com/ladmakhi81/learnup/shared/utils"
	"net/http"
)

type Handler struct {
	commentSvc     commentService.CommentService
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
}

func NewHandler(
	commentSvc commentService.CommentService,
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
) *Handler {
	return &Handler{
		commentSvc:     commentSvc,
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
	}
}

// GetCommentsPageable godoc
//
//	@Summary	Get paginated list of comments
//	@Tags		comments
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.GetCommentPageItemDto}}
//	@Failure	500			{object}	types.ApiError
//	@Router		/comments/page [get]
//	@Security	BearerAuth
func (h Handler) GetCommentsPageable(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	comments, count, err := h.commentSvc.Fetch(page, pageSize)
	if err != nil {
		return nil, err
	}
	commentsRes := types.NewPaginationRes(
		dtores.MapGetCommentPageItemsDto(comments),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, commentsRes), nil
}
