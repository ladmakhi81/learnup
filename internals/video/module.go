package video

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/video/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	videoHandler *handler.Handler
	middleware   *middleware.Middleware
}

func NewModule(
	handler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		videoHandler: handler,
		middleware:   middleware,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	videosApi := api.Group("/videos")
	videosApi.Use(m.middleware.CheckAccessToken())
	videosApi.PATCH("/:video-id/verify", utils.JsonHandler(m.videoHandler.VerifyVideo))
}
