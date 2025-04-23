package video

import (
	"github.com/gin-gonic/gin"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	videoHandler "github.com/ladmakhi81/learnup/internals/video/handler"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	videoHandler   *videoHandler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	userSvc userService.UserSvc,
	videoSvc videoService.VideoService,
	validationSvc contracts.Validation,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		videoHandler: videoHandler.NewHandler(
			validationSvc,
			videoSvc,
			translationSvc,
			userSvc,
		),
		middleware:     middleware,
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	videosApi := api.Group("/videos")
	videosApi.Use(m.middleware.CheckAccessToken())
	videosApi.PATCH("/:video-id/verify", utils.JsonHandler(m.translationSvc, m.videoHandler.VerifyVideo))
}
