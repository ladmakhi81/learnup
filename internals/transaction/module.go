package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/transaction/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	handler    *handler.Handler
	middleware *middleware.Middleware
}

func NewModule(
	handler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		handler:    handler,
		middleware: middleware,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	transactionsApi := api.Group("/transactions")
	transactionsApi.Use(m.middleware.CheckAccessToken())
	transactionsApi.GET("/page", utils.JsonHandler(m.handler.GetTransactionsPage))
}
