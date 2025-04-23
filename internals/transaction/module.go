package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/transaction/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	handler        *handler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	handler *handler.Handler,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		handler:        handler,
		middleware:     middleware,
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	transactionsApi := api.Group("/transactions")
	transactionsApi.Use(m.middleware.CheckAccessToken())
	transactionsApi.GET("/page", utils.JsonHandler(m.translationSvc, m.handler.GetTransactionsPage))
}
