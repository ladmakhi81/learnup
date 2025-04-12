package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/video/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/video/dto/res"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/pkg/validation"
	"github.com/ladmakhi81/learnup/types"
	"net/http"
)

type Handler struct {
	validationSvc  validation.Validation
	videoSvc       videoService.VideoService
	translationSvc translations.Translator
}

func NewHandler(
	validationSvc validation.Validation,
	videoSvc videoService.VideoService,
	translationSvc translations.Translator,
) *Handler {
	return &Handler{
		validationSvc:  validationSvc,
		videoSvc:       videoSvc,
		translationSvc: translationSvc,
	}
}

// AddNewVideoToCourse godoc
//
//	@Summary	Add a new video to a course
//	@Tags		videos
//	@Accept		json
//	@Produce	json
//	@Param		video	body		dtoreq.AddVideoToCourse	true	" "
//	@Success	201		{object}	types.ApiResponse{data=dtores.CreateCourseRes}
//	@Failure	400		{object}	types.ApiError
//	@Failure	409		{object}	types.ApiError
//	@Failure	500		{object}	types.ApiError
//	@Router		/videos/ [post]
func (h Handler) AddNewVideoToCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &dtoreq.AddVideoToCourse{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(h.translationSvc.Translate("common.errors.invalid_request_body"))
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	video, videoErr := h.videoSvc.AddVideo(dto)
	if videoErr != nil {
		return nil, videoErr
	}
	videoRes := dtores.NewCreateCourseRes(video.ID, video.URL, video.CourseId)
	return types.NewApiResponse(http.StatusCreated, videoRes), nil
}
