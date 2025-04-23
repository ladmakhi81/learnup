package video

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/video/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	videoHandler   *handler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	handler *handler.Handler,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		videoHandler:   handler,
		middleware:     middleware,
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	videosApi := api.Group("/videos")
	videosApi.Use(m.middleware.CheckAccessToken())
	videosApi.PATCH("/:video-id/verify", utils.JsonHandler(m.translationSvc, m.videoHandler.VerifyVideo))
}
