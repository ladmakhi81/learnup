package handler

import (
	"github.com/gin-gonic/gin"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	validationSvc  contracts.Validation
	videoSvc       videoService.VideoService
	translationSvc contracts.Translator
	userSvc        userService.UserSvc
}

func NewHandler(
	validationSvc contracts.Validation,
	videoSvc videoService.VideoService,
	translationSvc contracts.Translator,
	userSvc userService.UserSvc,
) *Handler {
	return &Handler{
		validationSvc:  validationSvc,
		videoSvc:       videoSvc,
		translationSvc: translationSvc,
		userSvc:        userSvc,
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
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	videoId, videoIdErr := utils.ToUint(ctx.Param("video-id"))
	if videoIdErr != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("video.errors.invalid_id"),
		)
	}
	if err := h.videoSvc.Verify(user, videoId); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, nil), nil
}
