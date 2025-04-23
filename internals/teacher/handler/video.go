package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	"github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/types"
	"net/http"
)

type VideoHandler struct {
	videoSvc       service.TeacherVideoService
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
}

func NewVideoHandler(
	videoSvc service.TeacherVideoService,
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
) *VideoHandler {
	return &VideoHandler{
		videoSvc:       videoSvc,
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
	}
}

// AddVideoToCourse godoc
//
//	@Summary	Add a new video to a course by teacher
//	@Tags		teacher
//	@Accept		json
//	@Produce	json
//	@Param		video	body		dtoreq.AddVideoToCourseReqDto	true	" "
//	@Success	201		{object}	types.ApiResponse{data=dtores.AddVideoToCourseResDto}
//	@Failure	400		{object}	types.ApiError
//	@Failure	401		{object}	types.ApiError
//	@Failure	409		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Router		/teacher/video [post]
//
//	@Security	BearerAuth
func (h VideoHandler) AddVideoToCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &dtoreq.AddVideoToCourseReqDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	video, err := h.videoSvc.AddVideo(*dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, dtores.NewAddVideoToCourseResDto(video)), nil
}
