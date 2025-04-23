package handler

import (
	"github.com/gin-gonic/gin"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	"github.com/ladmakhi81/learnup/internals/teacher/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type CommentHandler struct {
	teacherCommentSvc service.TeacherCommentService
	translationSvc    contracts.Translator
	userSvc           userService.UserSvc
}

func NewCommentHandler(
	teacherCommentSvc service.TeacherCommentService,
	translationSvc contracts.Translator,
	userSvc userService.UserSvc,
) *CommentHandler {
	return &CommentHandler{
		teacherCommentSvc: teacherCommentSvc,
		translationSvc:    translationSvc,
		userSvc:           userSvc,
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
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]dtores.GetCommentPageableItemDto}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/teacher/comments/{course-id} [get]
//	@Security	BearerAuth
func (h CommentHandler) GetPageableCommentByCourseId(ctx *gin.Context) (*types.ApiResponse, error) {
	courseId, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	comments, count, err := h.teacherCommentSvc.GetPageableCommentByCourseId(
		user,
		courseId,
		page,
		pageSize,
	)
	if err != nil {
		return nil, err
	}
	commentsRes := types.NewPaginationRes(
		dtores.MapGetCommentPageableItemsDto(comments),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, commentsRes), nil
}
