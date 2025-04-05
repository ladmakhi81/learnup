package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/auth/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type AuthModule struct {
	userAuthHandler *handler.UserAuthHandler
}

func NewAuthModule(userAuthHandler *handler.UserAuthHandler) *AuthModule {
	return &AuthModule{
		userAuthHandler: userAuthHandler,
	}
}

func (m AuthModule) Register(api *gin.RouterGroup) {
	authApi := api.Group("/auth")

	authApi.POST("/login", utils.JsonHandler(m.userAuthHandler.Login))
}
