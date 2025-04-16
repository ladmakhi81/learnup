package handler

import (
	"github.com/gin-gonic/gin"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"net/http"
	"strconv"
)

type Handler struct {
	validationSvc  contracts.Validation
	videoSvc       videoService.VideoService
	translationSvc contracts.Translator
}

func NewHandler(
	validationSvc contracts.Validation,
	videoSvc videoService.VideoService,
	translationSvc contracts.Translator,
) *Handler {
	return &Handler{
		validationSvc:  validationSvc,
		videoSvc:       videoSvc,
		translationSvc: translationSvc,
	}
}

// VerifyVideo godoc
//
//	@Summary	Verify a video
//	@Tags		videos
//	@Accept		json
//	@Produce	json
//	@Param		video-id	path		int	true " "
//	@Success	200			{object}	types.ApiResponse
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/videos/{video-id}/verify [patch]
//	@Security	BearerAuth
func (h Handler) VerifyVideo(ctx *gin.Context) (*types.ApiResponse, error) {
	authContext, _ := ctx.Get("AUTH")
	videoIdParam := ctx.Param("video-id")
	videoId, videoIdErr := strconv.Atoi(videoIdParam)
	if videoIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("video.errors.invalid_id"),
		)
	}
	if err := h.videoSvc.Verify(authContext, uint(videoId)); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}
