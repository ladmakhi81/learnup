package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	"github.com/ladmakhi81/learnup/internals/teacher/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/types"
	"github.com/ladmakhi81/learnup/shared/utils"
	"net/http"
)

type CourseHandler struct {
	courseSvc      service.TeacherCourseService
	validationSvc  contracts.Validation
	translationSvc contracts.Translator
	userSvc        userService.UserSvc
}

func NewCourseHandler(
	courseSvc service.TeacherCourseService,
	validationSvc contracts.Validation,
	translationSvc contracts.Translator,
	userSvc userService.UserSvc,
) *CourseHandler {
	return &CourseHandler{
		courseSvc:      courseSvc,
		validationSvc:  validationSvc,
		translationSvc: translationSvc,
		userSvc:        userSvc,
	}
}

// CreateCourse godoc
//
//	@Summary	Create a new course by teacher
//	@Tags		teacher
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtoreq.CreateCourseReqDto	true	" "
//	@Success	201		{object}	types.ApiResponse{data=dtores.CreateCourseResDto}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	409		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Router		/teacher/course [post]
//
//	@Security	BearerAuth
func (h CourseHandler) CreateCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &dtoreq.CreateCourseReqDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	teacher, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	course, err := h.courseSvc.Create(teacher, *dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, dtores.NewCreateCourseResDto(course)), nil
}

// FetchCourses godoc
//
//	@Summary	Get teacher's courses
//	@Tags		teacher
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page number"	default(0)
//	@Param		pageSize	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.FetchCourseItemDto}}
//	@Failure	401			{object}	types.ApiResponse
//	@Failure	404			{object}	types.ApiResponse
//	@Failure	500			{object}	types.ApiResponse
//	@Router		/teacher/courses [get]
//	@Security	BearerAuth
func (h CourseHandler) FetchCourses(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Param("page"), ctx.Param("pageSize"))
	teacher, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	courses, count, err := h.courseSvc.FetchByTeacherId(teacher, page, pageSize)
	if err != nil {
		return nil, err
	}
	coursesRes := types.NewPaginationRes(
		dtores.MapFetchCourseItemsDto(courses),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, coursesRes), nil
}
