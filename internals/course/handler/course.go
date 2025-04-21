package handler

import (
	"github.com/gin-gonic/gin"
	commentDtoReq "github.com/ladmakhi81/learnup/internals/comment/dto/req"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	courseDtoReq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	questionDtoRes "github.com/ladmakhi81/learnup/internals/course/dto/res"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	likeDtoReq "github.com/ladmakhi81/learnup/internals/like/dto/req"
	likeService "github.com/ladmakhi81/learnup/internals/like/service"
	questionDtoReq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	questionService "github.com/ladmakhi81/learnup/internals/question/service"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
	"strconv"
)

type Handler struct {
	courseSvc     courseService.CourseService
	validationSvc contracts.Validation
	translateSvc  contracts.Translator
	videoSvc      videoService.VideoService
	likeSvc       likeService.LikeService
	commentSvc    commentService.CommentService
	questionSvc   questionService.QuestionService
}

func NewHandler(
	courseSvc courseService.CourseService,
	validationSvc contracts.Validation,
	translateSvc contracts.Translator,
	videosSvc videoService.VideoService,
	likeSvc likeService.LikeService,
	commentSvc commentService.CommentService,
	questionSvc questionService.QuestionService,
) *Handler {
	return &Handler{
		courseSvc:     courseSvc,
		validationSvc: validationSvc,
		translateSvc:  translateSvc,
		videoSvc:      videosSvc,
		likeSvc:       likeSvc,
		commentSvc:    commentSvc,
		questionSvc:   questionSvc,
	}
}

// CreateCourse godoc
//
//	@Summary	Create a new course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		requestBody	body		courseDtoReq.CreateCourseReq	true	" "
//	@Success	201			{object}	types.ApiResponse{data=questionDtoRes.CreateCourseRes}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	409			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses [post]
//	@Security	BearerAuth
func (h Handler) CreateCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &courseDtoReq.CreateCourseReq{}
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	authContext, _ := ctx.Get("AUTH")
	course, courseErr := h.courseSvc.Create(authContext, *dto)
	if courseErr != nil {
		return nil, courseErr
	}
	courseRes := questionDtoRes.CreateCourseRes{
		ID:                          course.ID,
		Fee:                         course.Fee,
		Price:                       course.Price,
		VerifiedByID:                course.VerifiedByID,
		VerifiedDate:                course.VerifiedDate,
		TeacherID:                   course.TeacherID,
		ThumbnailImage:              course.ThumbnailImage,
		Tags:                        course.Tags,
		Status:                      course.Status,
		Prerequisite:                course.Prerequisite,
		MaxDiscountAmount:           course.MaxDiscountAmount,
		Level:                       course.Level,
		IsVerifiedByAdmin:           course.IsVerifiedByAdmin,
		IntroductionVideo:           course.IntroductionVideo,
		Image:                       course.Image,
		IsPublished:                 course.IsPublished,
		DiscountFeeAmountPercentage: course.DiscountFeeAmountPercentage,
		Description:                 course.Description,
		CommentAccessMode:           course.CommentAccessMode,
		CanHaveDiscount:             course.CanHaveDiscount,
		AbilityToAddComment:         course.AbilityToAddComment,
		Name:                        course.Name,
		CategoryID:                  course.CategoryID,
	}
	return types.NewApiResponse(http.StatusCreated, courseRes), nil
}

// GetCourses godoc
//
//	@Summary	Get list of paginated courses
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]questionDtoRes.GetCoursesRes}}
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
	courses, coursesErr := h.courseSvc.GetCourses(page, pageSize)
	if coursesErr != nil {
		return nil, coursesErr
	}
	coursesCount, coursesCountErr := h.courseSvc.GetCoursesCount()
	if coursesCountErr != nil {
		return nil, coursesCountErr
	}
	mappedCourses := questionDtoRes.NewGetCoursesRes(courses)
	paginationRes := types.NewPaginationRes(
		mappedCourses,
		page,
		utils.CalculatePaginationTotalPage(coursesCount, pageSize),
		coursesCount,
	)
	return types.NewApiResponse(http.StatusOK, paginationRes), nil
}

// GetVideosByCourseID godoc
//
//	@Summary	Get Videos by Course ID
//	@Tags		courses
//	@Param		course-id	path		int	true	"Course ID"
//	@Success	200			{object}	questionDtoRes.GetVideosByCourseIDRes
//	@Failure	400			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/videos [get]
//
//	@Security	BearerAuth
func (h Handler) GetVideosByCourseID(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIDParam := ctx.Param("course-id")
	courseID, courseIDErr := strconv.Atoi(courseIDParam)
	if courseIDErr != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	videos, videosErr := h.videoSvc.FindVideosByCourseID(uint(courseID))
	if videosErr != nil {
		return nil, videosErr
	}
	videosRes := questionDtoRes.NewGetVideosByCourseIDRes(videos, uint(courseID))
	return types.NewApiResponse(http.StatusOK, videosRes), nil
}

// GetCourseById godoc
//
//	@Summary	Get Course by ID
//	@Tags		courses
//	@Param		course-id	path		int	true	"Course ID"
//	@Success	200			{object}	questionDtoRes.GetCourseByIdRes
//	@Failure	400			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id} [get]
//
//	@Security	BearerAuth
func (h Handler) GetCourseById(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIDParam := ctx.Param("course-id")
	courseID, courseIDErr := strconv.Atoi(courseIDParam)
	if courseIDErr != nil {
		return nil, types.NewBadRequestError(h.translateSvc.Translate("course.errors.invalid_course_id"))
	}
	course, courseErr := h.courseSvc.FindDetailById(uint(courseID))
	if courseErr != nil {
		return nil, courseErr
	}
	courseRes := questionDtoRes.NewGetCourseByIdRes(course)
	return types.NewApiResponse(http.StatusOK, courseRes), nil
}

// VerifyCourse godoc
//
//	@Summary	Verify a course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int								true	"Course ID"
//	@Param		request		body		courseDtoReq.VerifyCourseReq	true	" "
//	@Success	200			{object}	types.ApiResponse
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/verify [patch]
//
//	@Security	BearerAuth
func (h Handler) VerifyCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIdParam := ctx.Param("course-id")
	courseId, courseIdErr := strconv.Atoi(courseIdParam)
	if courseIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	dto := &courseDtoReq.VerifyCourseReq{
		ID: uint(courseId),
	}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	authContext, _ := ctx.Get("AUTH")
	if err := h.courseSvc.VerifyCourse(authContext, *dto); err != nil {
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
//	@Param		like		body		likeDtoReq.CreateLikeReq	true	" "
//	@Success	201			{object}	types.ApiResponse
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/like [post]
//
//	@Security	BearerAuth
func (h Handler) Like(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIDParam := ctx.Param("course-id")
	courseID, courseIDErr := strconv.Atoi(courseIDParam)
	if courseIDErr != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	authContext, _ := ctx.Get("AUTH")
	dto := &likeDtoReq.CreateLikeReq{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	dto.CourseID = uint(courseID)
	_, likeErr := h.likeSvc.Create(authContext, *dto)
	if likeErr != nil {
		return nil, likeErr
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
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]questionDtoRes.GetLikesPageableItem}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/likes [get]
//	@Security	BearerAuth
func (h Handler) FetchLikes(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIDParam := ctx.Param("course-id")
	courseID, courseIDErr := strconv.Atoi(courseIDParam)
	if courseIDErr != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	likes, likesErr := h.likeSvc.FetchByCourseID(page, pageSize, uint(courseID))
	if likesErr != nil {
		return nil, likesErr
	}
	likeCount, likeCountErr := h.likeSvc.FetchCountByCourseID(uint(courseID))
	if likeCountErr != nil {
		return nil, likeCountErr
	}
	likesRes := types.NewPaginationRes(
		questionDtoRes.MappedGetLikesPageableItem(likes),
		page,
		utils.CalculatePaginationTotalPage(likeCount, pageSize),
		likeCount,
	)
	return types.NewApiResponse(http.StatusOK, likesRes), nil
}

// CreateComment godoc
//
//	@Summary	Create a new comment
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		course-id	path		int								true	"Course ID"
//	@Param		request		body		commentDtoReq.CreateCommentReq	true	" "
//	@Success	201			{object}	types.ApiResponse{data=questionDtoRes.CreateCommentRes}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/comment [post]
//	@Security	BearerAuth
func (h Handler) CreateComment(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIDParam := ctx.Param("course-id")
	courseID, courseIDErr := strconv.Atoi(courseIDParam)
	if courseIDErr != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	dto := &commentDtoReq.CreateCommentReq{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	dto.CourseId = uint(courseID)
	authContext, _ := ctx.Get("AUTH")
	comment, commentErr := h.commentSvc.Create(authContext, *dto)
	if commentErr != nil {
		return nil, commentErr
	}
	commentRes := questionDtoRes.NewCreateCommentRes(comment)
	return types.NewApiResponse(http.StatusCreated, commentRes), nil
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
	commentIdParam := ctx.Param("comment-id")
	commentId, commentIdErr := strconv.Atoi(commentIdParam)
	if commentIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("comment.errors.invalid_id"),
		)
	}
	if err := h.commentSvc.Delete(uint(commentId)); err != nil {
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
//	@Param		question	body		questionDtoReq.CreateQuestionReq	true	" "
//	@Success	201			{object}	types.ApiResponse{data=questionDtoRes.CreateQuestionRes}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/question [post]
//	@Security	BearerAuth
func (h Handler) CreateQuestion(ctx *gin.Context) (*types.ApiResponse, error) {
	authContext, _ := ctx.Get("AUTH")
	senderClaim, _ := authContext.(*types.TokenClaim)
	courseIDParam := ctx.Param("course-id")
	courseID, courseIDErr := strconv.Atoi(courseIDParam)
	if courseIDErr != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	dto := &questionDtoReq.CreateQuestionReq{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	dto.CourseID = uint(courseID)
	dto.UserID = senderClaim.UserID
	question, questionErr := h.questionSvc.Create(*dto)
	if questionErr != nil {
		return nil, questionErr
	}
	questionRes := questionDtoRes.NewCreateQuestionRes(question)
	return types.NewApiResponse(http.StatusCreated, questionRes), nil
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
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{rows=[]questionDtoRes.GetQuestionItemRes}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses/{course-id}/questions [get]
//	@Security	BearerAuth
func (h Handler) GetQuestions(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIDParam := ctx.Param("course-id")
	parsedCourseId, parsedCourseIdErr := strconv.Atoi(courseIDParam)
	if parsedCourseIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("course.errors.invalid_course_id"),
		)
	}
	courseID := uint(parsedCourseId)
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	questions, questionsErr := h.questionSvc.GetPageable(&courseID, page, pageSize)
	if questionsErr != nil {
		return nil, questionsErr
	}
	questionCount, questionCountErr := h.questionSvc.GetCount(&courseID)
	if questionCountErr != nil {
		return nil, questionCountErr
	}
	questionRes := types.NewPaginationRes(
		questionDtoRes.MapGetQuestionItemRes(questions),
		page,
		utils.CalculatePaginationTotalPage(questionCount, pageSize),
		questionCount,
	)
	return types.NewApiResponse(http.StatusOK, questionRes), nil
}
