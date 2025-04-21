package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	"github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type CourseHandler struct {
	courseSvc      service.TeacherCourseService
	validationSvc  contracts.Validation
	translationSvc contracts.Translator
}

func NewCourseHandler(
	courseSvc service.TeacherCourseService,
	validationSvc contracts.Validation,
	translationSvc contracts.Translator,
) *CourseHandler {
	return &CourseHandler{
		courseSvc:      courseSvc,
		validationSvc:  validationSvc,
		translationSvc: translationSvc,
	}
}

// CreateCourse godoc
//
//	@Summary	Create a new course by teacher
//	@Tags		teacher
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtoreq.CreateCourseReq	true	" "
//	@Success	201		{object}	types.ApiResponse{data=dtores.CreateCourseRes}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	409		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Router		/teacher/course [post]
//
//	@Security	BearerAuth
func (h CourseHandler) CreateCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &dtoreq.CreateCourseReq{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
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
		CreatedAt: course.CreatedAt,
		ID:        course.ID,
		UpdatedAt: course.UpdatedAt,
	}
	return types.NewApiResponse(http.StatusCreated, courseRes), nil
}

// FetchCourses godoc
//
//	@Summary	Get teacher's courses
//	@Tags		teacher
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.FetchCourseItemRes}}
//	@Failure	401			{object}	types.ApiResponse
//	@Failure	404			{object}	types.ApiResponse
//	@Failure	500			{object}	types.ApiResponse
//	@Router		/teacher/courses [get]
//	@Security	BearerAuth
func (h CourseHandler) FetchCourses(ctx *gin.Context) (*types.ApiResponse, error) {
	authContext, _ := ctx.Get("AUTH")
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Param("page"), ctx.Param("pageSize"))
	courses, count, coursesErr := h.courseSvc.FetchByTeacherId(authContext, page, pageSize)
	if coursesErr != nil {
		return nil, coursesErr
	}
	coursesRes := types.NewPaginationRes(
		dtores.MapCoursesToFetchCourseItemRes(courses),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, coursesRes), nil
}
