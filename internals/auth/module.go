package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/auth/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	userAuthHandler *handler.UserAuthHandler
}

func NewModule(userAuthHandler *handler.UserAuthHandler) *Module {
	return &Module{
		userAuthHandler: userAuthHandler,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	authApi := api.Group("/auth")

	authApi.POST("/login", utils.JsonHandler(m.userAuthHandler.Login))
}
