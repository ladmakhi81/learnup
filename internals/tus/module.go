package tus

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/tus/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	hookHandler    *handler.TusHookHandler
	translationSvc contracts.Translator
}

func NewModule(hookHandler *handler.TusHookHandler, translationSvc contracts.Translator) *Module {
	return &Module{
		hookHandler:    hookHandler,
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	tusHookApi := api.Group("/tus-hooks")

	tusHookApi.POST("/videos", utils.JsonHandler(m.translationSvc, m.hookHandler.VideoWebhook))
}
