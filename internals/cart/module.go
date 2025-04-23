package cart

import (
	"github.com/gin-gonic/gin"
	cartHandler "github.com/ladmakhi81/learnup/internals/cart/handler"
	cartService "github.com/ladmakhi81/learnup/internals/cart/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	handler        *cartHandler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	userSvc userService.UserSvc,
	cartSvc cartService.CartService,
	validationSvc contracts.Validation,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		handler:        cartHandler.NewHandler(translationSvc, validationSvc, cartSvc, userSvc),
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
