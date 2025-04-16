package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	"github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"net/http"
)

type Handler struct {
	courseSvc      service.TeacherCourseService
	validationSvc  contracts.Validation
	translationSvc contracts.Translator
}

func NewHandler(
	courseSvc service.TeacherCourseService,
	validationSvc contracts.Validation,
	translationSvc contracts.Translator,
) *Handler {
	return &Handler{
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
func (h Handler) CreateCourse(ctx *gin.Context) (*types.ApiResponse, error) {
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
