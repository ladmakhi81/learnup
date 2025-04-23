package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/auth/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	userAuthHandler *handler.Handler
	translationSvc  contracts.Translator
}

func NewModule(translationSvc contracts.Translator, userAuthHandler *handler.Handler) *Module {
	return &Module{
		userAuthHandler: userAuthHandler,
		translationSvc:  translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	authApi := api.Group("/auth")

	authApi.POST("/login", utils.JsonHandler(m.translationSvc, m.userAuthHandler.Login))
}
