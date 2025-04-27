package handler

import (
	"github.com/gin-gonic/gin"
	commentDtoReq "github.com/ladmakhi81/learnup/internals/comment/dto/req"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	courseDtoReq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	courseDtoRes "github.com/ladmakhi81/learnup/internals/course/dto/res"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	forumService "github.com/ladmakhi81/learnup/internals/forum/service"
	likeDtoReq "github.com/ladmakhi81/learnup/internals/like/dto/req"
	likeService "github.com/ladmakhi81/learnup/internals/like/service"
	questionDtoReq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	questionService "github.com/ladmakhi81/learnup/internals/question/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/types"
	"github.com/ladmakhi81/learnup/shared/utils"
	"net/http"
)

type Handler struct {
	courseSvc     courseService.CourseService
	validationSvc contracts.Validation
	translateSvc  contracts.Translator
	videoSvc      videoService.VideoService
	likeSvc       likeService.LikeService
	commentSvc    commentService.CommentService
	questionSvc   questionService.QuestionService
	userSvc       userService.UserSvc
	forumSvc      forumService.ForumService
}

func NewHandler(
	courseSvc courseService.CourseService,
	validationSvc contracts.Validation,
	translateSvc contracts.Translator,
	videosSvc videoService.VideoService,
	likeSvc likeService.LikeService,
	commentSvc commentService.CommentService,
	questionSvc questionService.QuestionService,
	userSvc userService.UserSvc,
	forumSvc forumService.ForumService,
) *Handler {
	return &Handler{
		courseSvc:     courseSvc,
		validationSvc: validationSvc,
		translateSvc:  translateSvc,
		videoSvc:      videosSvc,
		likeSvc:       likeSvc,
		commentSvc:    commentSvc,
		questionSvc:   questionSvc,
		userSvc:       userSvc,
		forumSvc:      forumSvc,
	}
}

// CreateCourse godoc
//
//	@Summary	Create a new course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		requestBody	body		courseDtoReq.CreateCourseReqDto	true	" "
//	@Success	201			{object}	types.ApiResponse{data=courseDtoRes.CreateCourseResDto}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	409			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses [post]
//	@Security	BearerAuth
func (h Handler) CreateCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &courseDtoReq.CreateCourseReqDto{}
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	course, err := h.courseSvc.Create(user, *dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, courseDtoRes.NewCreateCourseResDto(course)), nil
}

// GetCourses godoc
//
//	@Summary	Get list of paginated courses
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]courseDtoRes.GetPageableCourseItemDto}}
//	@Failure	401			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/page [get]
//
//	@Security	BearerAuth
func (h Handler) GetCourses(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	courses, count, err := h.courseSvc.GetCourses(page, pageSize)
	if err != nil {
		return nil, err
	}
	paginationRes := types.NewPaginationRes(
		courseDtoRes.MapGetPageableCourseItemsDto(courses),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, paginationRes), nil
}

// GetVideosByCourseID godoc
//
//	@Summary	Get Videos by Course ID
//	@Tags		courses
//	@Param		course-id	path		int	true	"Course ID"
//	@Success	200			{object}	[]courseDtoRes.GetVideoByCourseItemDto
//	@Failure	400			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/videos [get]
//
//	@Security	BearerAuth
func (h Handler) GetVideosByCourseID(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	videos, err := h.videoSvc.FindVideosByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, courseDtoRes.MapGetVideoByCourseItemsDto(videos)), nil
}

// GetCourseById godoc
//
//	@Summary	Get Course by ID
//	@Tags		courses
//	@Param		course-id	path		int	true	"Course ID"
//	@Success	200			{object}	courseDtoRes.GetCourseByItemDto
//	@Failure	400			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id} [get]
//
//	@Security	BearerAuth
func (h Handler) GetCourseById(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	course, err := h.courseSvc.FindDetailById(courseID)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, courseDtoRes.NewGetCourseByItemDto(course)), nil
}

// VerifyCourse godoc
//
//	@Summary	Verify a course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int								true	"Course ID"
//	@Param		request		body		courseDtoReq.VerifyCourseReqDto	true	" "
//	@Success	200			{object}	types.ApiResponse
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/verify [patch]
//
//	@Security	BearerAuth
func (h Handler) VerifyCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	dto := &courseDtoReq.VerifyCourseReqDto{
		ID: courseID,
	}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.courseSvc.VerifyCourse(user, *dto); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

// Like godoc
//
//	@Summary	Like a course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int							true	"Course ID"
//	@Param		like		body		likeDtoReq.CreateLikeReqDto	true	" "
//	@Success	201			{object}	types.ApiResponse
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/like [post]
//
//	@Security	BearerAuth
func (h Handler) Like(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	dto := &likeDtoReq.CreateLikeReqDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	dto.CourseID = courseID
	if _, err := h.likeSvc.Create(user, *dto); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, nil), nil
}

// FetchLikes godoc
//
//	@Summary	Get paginated likes by course ID
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int	true	"Course ID"
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]courseDtoRes.GetLikesPageableItemDto}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/likes [get]
//	@Security	BearerAuth
func (h Handler) FetchLikes(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	likes, count, err := h.likeSvc.FetchByCourseID(page, pageSize, courseID)
	if err != nil {
		return nil, err
	}
	likesRes := types.NewPaginationRes(
		courseDtoRes.MapGetLikesPageableItemsDto(likes),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, likesRes), nil
}

// CreateComment godoc
//
//	@Summary	Create a new comment
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int									true	"Course ID"
//	@Param		request		body		commentDtoReq.CreateCommentReqDto	true	" "
//	@Success	201			{object}	types.ApiResponse{data=courseDtoRes.CreateCommentResDto}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/comment [post]
//	@Security	BearerAuth
func (h Handler) CreateComment(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	dto := &commentDtoReq.CreateCommentReqDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	dto.CourseId = courseID
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	comment, err := h.commentSvc.Create(user, *dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, courseDtoRes.NewCreateCommentResDto(comment)), nil
}

// DeleteComment godoc
//
//	@Summary	Delete a comment
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		comment-id	path		int	true	"Comment ID"
//	@Param		course-id	path		int	true	"Course ID"
//	@Success	200			{object}	types.ApiResponse
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/comments/{comment-id} [delete]
//	@Security	BearerAuth
func (h Handler) DeleteComment(ctx *gin.Context) (*types.ApiResponse, error) {
	commentID, err := utils.ToUint(ctx.Param("comment-id"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("comment.errors.invalid_id"),
		)
	}
	if err := h.commentSvc.Delete(commentID); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}

// CreateQuestion godoc
//
//	@Summary	Create a new question for a course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int									true	"Course ID"
//	@Param		question	body		questionDtoReq.CreateQuestionReqDto	true	" "
//	@Success	201			{object}	types.ApiResponse{data=courseDtoRes.CreateQuestionResDto}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/question [post]
//	@Security	BearerAuth
func (h Handler) CreateQuestion(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	dto := &questionDtoReq.CreateQuestionReqDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	dto.CourseID = courseID
	question, err := h.questionSvc.Create(user, *dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, courseDtoRes.NewCreateQuestionResDto(question)), nil
}

// GetQuestions godoc
//
//	@Summary	Retrieve paginated questions for a specific course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int	true	"Course ID"
//	@Param		page		query		int	false	"Page number"				default(0)
//	@Param		pageSize	query		int	false	"Number of items per page"	default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]courseDtoRes.GetQuestionItemDto}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/questions [get]
//	@Security	BearerAuth
func (h Handler) GetQuestions(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	questions, count, err := h.questionSvc.GetPageable(&courseID, page, pageSize)
	if err != nil {
		return nil, err
	}
	questionRes := types.NewPaginationRes(
		courseDtoRes.MapGetQuestionItemsDto(questions),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, questionRes), nil
}

// GetForumByCourseID godoc
// @Summary	Get Forum by Course ID
// @Tags		courses
// @Accept		json
// @Produce	json
// @Param		course-id	path		int	true	"Course ID"
// @Success	200			{object}	types.ApiResponse{data=courseDtoRes.GetForumByCourseIDDto}
// @Failure	400			{object}	types.ApiError
// @Failure	500			{object}	types.ApiError
// @Router		/courses/{course-id}/forum [get]
// @Security	BearerAuth
func (h Handler) GetForumByCourseID(ctx *gin.Context) (*types.ApiResponse, error) {
	courseID, err := utils.ToUint(ctx.Param("course-id"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	forum, err := h.forumSvc.GetForumByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, courseDtoRes.MapGetForumByCourseIDDto(forum)), nil
}
