package tus

import (
	"github.com/gin-gonic/gin"
	tusHandler "github.com/ladmakhi81/learnup/internals/tus/handler"
	tusHookService "github.com/ladmakhi81/learnup/internals/tus/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	hookHandler    *tusHandler.TusHookHandler
	translationSvc contracts.Translator
}

func NewModule(tusHookSvc tusHookService.TusService, translationSvc contracts.Translator) *Module {
	return &Module{
		hookHandler:    tusHandler.NewTusHookHandler(tusHookSvc),
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	tusHookApi := api.Group("/tus-hooks")

	tusHookApi.POST("/videos", utils.JsonHandler(m.translationSvc, m.hookHandler.VideoWebhook))
}
