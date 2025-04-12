package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/course/dto/res"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/pkg/validation"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
	"strconv"
)

type Handler struct {
	courseSvc     courseService.CourseService
	validationSvc validation.Validation
	translateSvc  translations.Translator
	videoSvc      videoService.VideoService
}

func NewHandler(
	courseSvc courseService.CourseService,
	validationSvc validation.Validation,
	translateSvc translations.Translator,
	videosSvc videoService.VideoService,
) *Handler {
	return &Handler{
		courseSvc:     courseSvc,
		validationSvc: validationSvc,
		translateSvc:  translateSvc,
		videoSvc:      videosSvc,
	}
}

// CreateCourse godoc
//
//	@Summary	Create a new course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		requestBody	body		dtoreq.CreateCourseReq	true	" "
//	@Success	201			{object}	types.ApiResponse{data=dtores.CreateCourseRes}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	409			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/courses [post]
//	@Security	BearerAuth
func (h Handler) CreateCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &dtoreq.CreateCourseReq{}
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
	courseRes := dtores.CreateCourseRes{
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
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.GetCoursesRes}}
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
	mappedCourses := dtores.NewGetCoursesRes(courses)
	paginationRes := types.NewPaginationRes(
		mappedCourses,
		page,
		utils.CalculatePaginationTotalPage(coursesCount),
		coursesCount,
	)
	return types.NewApiResponse(http.StatusOK, paginationRes), nil
}

// GetVideosByCourseID godoc
//
//	@Summary	Get Videos by Course ID
//	@Tags		courses
//	@Param		course-id	path		int	true	"Course ID"
//	@Success	200			{object}	dtores.GetVideosByCourseIDRes
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
	videosRes := dtores.NewGetVideosByCourseIDRes(videos, uint(courseID))
	return types.NewApiResponse(http.StatusOK, videosRes), nil
}

// GetCourseById godoc
//
//	@Summary	Get Course by ID
//	@Tags		courses
//	@Param		course-id	path		int	true	"Course ID"
//	@Success	200			{object}	dtores.GetCourseByIdRes
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
	courseRes := dtores.NewGetCourseByIdRes(course)
	return types.NewApiResponse(http.StatusOK, courseRes), nil
}
