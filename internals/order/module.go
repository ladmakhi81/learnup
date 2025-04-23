package order

import (
	"github.com/gin-gonic/gin"
	orderHandler "github.com/ladmakhi81/learnup/internals/order/handler"
	orderService "github.com/ladmakhi81/learnup/internals/order/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	handler        *orderHandler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	orderSvc orderService.OrderService,
	validationSvc contracts.Validation,
	userSvc userService.UserSvc,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		handler:        orderHandler.NewHandler(orderSvc, translationSvc, validationSvc, userSvc),
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
