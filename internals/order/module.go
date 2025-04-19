package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/order/handler"
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
	ordersApi := api.Group("/orders")
	ordersApi.Use(m.middleware.CheckAccessToken())

	ordersApi.POST("/", utils.JsonHandler(m.handler.CreateOrder))
	ordersApi.GET("/", utils.JsonHandler(m.handler.GetOrdersPage))
	ordersApi.GET("/:order-id", utils.JsonHandler(m.handler.GetOrderByID))
}
