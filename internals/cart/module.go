package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/cart/handler"
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
	cartsApi := api.Group("/carts")
	cartsApi.Use(m.middleware.CheckAccessToken())
	cartsApi.POST("/", utils.JsonHandler(m.translationSvc, m.handler.AddCart))
	cartsApi.GET("/", utils.JsonHandler(m.translationSvc, m.handler.GetCartsByUserID))
	cartsApi.DELETE("/:cart-id", utils.JsonHandler(m.translationSvc, m.handler.DeleteCartByID))
}
