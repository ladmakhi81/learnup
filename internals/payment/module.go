package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/payment/handler"
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
	paymentsApi := api.Group("/payments")

	paymentsApi.GET("/verify/zarinpal", utils.JsonHandler(m.handler.VerifyZarinpal))
	paymentsApi.GET("/verify/zibal", utils.JsonHandler(m.handler.VerifyZibal))
	paymentsApi.GET("/verify/stripe", utils.JsonHandler(m.handler.VerifyStripe))
	paymentsApi.GET("/page", m.middleware.CheckAccessToken(), utils.JsonHandler(m.handler.GetPayments))
}
