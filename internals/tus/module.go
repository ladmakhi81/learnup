package tus

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/tus/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	hookHandler *handler.TusHookHandler
}

func NewModule(hookHandler *handler.TusHookHandler) *Module {
	return &Module{
		hookHandler: hookHandler,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	tusHookApi := api.Group("/tus-hooks")

	tusHookApi.POST("/videos", utils.JsonHandler(m.hookHandler.VideoWebhook))
}
