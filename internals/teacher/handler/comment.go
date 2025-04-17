package handler

import (
	"github.com/gin-gonic/gin"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	"github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	teacherCommentSvc service.TeacherCommentService
	translationSvc    contracts.Translator
}

func NewCommentHandler(
	teacherCommentSvc service.TeacherCommentService,
	translationSvc contracts.Translator,
) *CommentHandler {
	return &CommentHandler{
		teacherCommentSvc: teacherCommentSvc,
		translationSvc:    translationSvc,
	}
}

// GetPageableCommentByCourseId godoc
//
//	@Summary	Get paginated comments by course ID
//	@Tags		teacher
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int	true	"Course ID"
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]dtores.GetCommentPageableItemRes}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/teacher/comments/{course-id} [get]
//	@Security	BearerAuth
func (h CommentHandler) GetPageableCommentByCourseId(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIdParam := ctx.Param("course-id")
	courseId, courseIdErr := strconv.Atoi(courseIdParam)
	if courseIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	authContext, _ := ctx.Get("AUTH")
	comments, commentsErr := h.teacherCommentSvc.GetPageableCommentByCourseId(
		authContext,
		uint(courseId),
		page,
		pageSize,
	)
	if commentsErr != nil {
		return nil, commentsErr
	}
	commentCount, commentCountErr := h.teacherCommentSvc.GetCommentCountByCourseId(
		authContext,
		uint(courseId),
	)
	if commentCountErr != nil {
		return nil, commentCountErr
	}
	commentsRes := types.NewPaginationRes(
		dtores.MappedGetCommentPageableItemsRes(comments),
		page,
		utils.CalculatePaginationTotalPage(commentCount),
		commentCount,
	)
	return types.NewApiResponse(http.StatusOK, commentsRes), nil
}
