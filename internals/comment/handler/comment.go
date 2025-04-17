package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/comment/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/comment/dto/res"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
	"strconv"
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

// CreateComment godoc
//
//	@Summary	Create a new comment
//	@Tags		comments
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtoreq.CreateCommentReq	true	" "
//	@Success	201		{object}	types.ApiResponse{data=dtores.CreateCommentRes}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	404		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Router		/comments [post]
//	@Security	BearerAuth
func (h Handler) CreateComment(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &dtoreq.CreateCommentReq{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	authContext, _ := ctx.Get("AUTH")
	comment, commentErr := h.commentSvc.Create(authContext, *dto)
	if commentErr != nil {
		return nil, commentErr
	}
	commentRes := dtores.NewCreateCommentRes(comment)
	return types.NewApiResponse(http.StatusCreated, commentRes), nil
}

// DeleteComment godoc
//
//	@Summary	Delete a comment
//	@Tags		comments
//	@Accept		json
//	@Produce	json
//	@Param		comment-id	path		int	true	"Comment ID"
//	@Success	200			{object}	types.ApiResponse
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/comments/{comment-id} [delete]
//	@Security	BearerAuth
func (h Handler) DeleteComment(ctx *gin.Context) (*types.ApiResponse, error) {
	commentIdParam := ctx.Param("comment-id")
	commentId, commentIdErr := strconv.Atoi(commentIdParam)
	if commentIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("comment.errors.invalid_id"),
		)
	}
	if err := h.commentSvc.Delete(uint(commentId)); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

// GetCommentsPageable godoc
//
//	@Summary	Get paginated list of comments
//	@Tags		comments
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.GetCommentPageItem}}
//	@Failure	500			{object}	types.ApiError
//	@Router		/comments/page [get]
//	@Security	BearerAuth
func (h Handler) GetCommentsPageable(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	comments, commentsErr := h.commentSvc.Fetch(page, pageSize)
	if commentsErr != nil {
		return nil, commentsErr
	}
	commentsCount, commentsCountErr := h.commentSvc.FetchCount()
	if commentsCountErr != nil {
		return nil, commentsCountErr
	}
	commentsRes := types.NewPaginationRes(
		dtores.NewGetCommentsPageableItem(comments),
		page,
		utils.CalculatePaginationTotalPage(commentsCount),
		commentsCount,
	)
	return types.NewApiResponse(http.StatusOK, commentsRes), nil
}
