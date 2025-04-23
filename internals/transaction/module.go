package transaction

import (
	"github.com/gin-gonic/gin"
	transactionHandler "github.com/ladmakhi81/learnup/internals/transaction/handler"
	transactionService "github.com/ladmakhi81/learnup/internals/transaction/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	handler        *transactionHandler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	transactionSvc transactionService.TransactionService,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		handler:        transactionHandler.NewHandler(transactionSvc),
		middleware:     middleware,
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	transactionsApi := api.Group("/transactions")
	transactionsApi.Use(m.middleware.CheckAccessToken())
	transactionsApi.GET("/page", utils.JsonHandler(m.translationSvc, m.handler.GetTransactionsPage))
}
