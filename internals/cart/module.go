package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/cart/handler"
	"github.com/ladmakhi81/learnup/internals/middleware"
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
	cartsApi := api.Group("/carts")
	cartsApi.Use(m.middleware.CheckAccessToken())
	cartsApi.POST("/", utils.JsonHandler(m.handler.AddCart))
	// TODO: move these endpoints under user profile endpoint
	cartsApi.GET("/", utils.JsonHandler(m.handler.GetCartsByUserID))
	cartsApi.DELETE("/:cart-id", utils.JsonHandler(m.handler.DeleteCartByID))
}
