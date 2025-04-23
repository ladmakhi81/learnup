package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/order/handler"
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
	ordersApi := api.Group("/orders")
	ordersApi.Use(m.middleware.CheckAccessToken())

	ordersApi.POST("/", utils.JsonHandler(m.translationSvc, m.handler.CreateOrder))
	ordersApi.GET("/", utils.JsonHandler(m.translationSvc, m.handler.GetOrdersPage))
	ordersApi.GET("/:order-id", utils.JsonHandler(m.translationSvc, m.handler.GetOrderByID))
}
