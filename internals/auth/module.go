package auth

import (
	"github.com/gin-gonic/gin"
	authHandler "github.com/ladmakhi81/learnup/internals/auth/handler"
	authService "github.com/ladmakhi81/learnup/internals/auth/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	authHandler    *authHandler.Handler
	translationSvc contracts.Translator
}

func NewModule(
	authSvc authService.AuthService,
	validationSvc contracts.Validation,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		authHandler:    authHandler.NewHandler(authSvc, validationSvc, translationSvc),
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	authApi := api.Group("/auth")
	authApi.POST("/login", utils.JsonHandler(m.translationSvc, m.authHandler.Login))
}
