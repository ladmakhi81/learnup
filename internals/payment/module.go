package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/payment/handler"
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
	paymentsApi := api.Group("/payments")

	paymentsApi.GET("/verify/zarinpal", utils.JsonHandler(m.translationSvc, m.handler.VerifyZarinpal))
	paymentsApi.GET("/verify/zibal", utils.JsonHandler(m.translationSvc, m.handler.VerifyZibal))
	paymentsApi.GET("/verify/stripe", utils.JsonHandler(m.translationSvc, m.handler.VerifyStripe))
	paymentsApi.GET("/page", m.middleware.CheckAccessToken(), utils.JsonHandler(m.translationSvc, m.handler.GetPayments))
}
