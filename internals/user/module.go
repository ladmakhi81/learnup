package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/user/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	userAdminHandler *handler.Handler
	middleware       *middleware.Middleware
}

func NewModule(
	userAdminHandler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		userAdminHandler: userAdminHandler,
		middleware:       middleware,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	usersApi := api.Group("/users")

	usersApi.POST("/basic", utils.JsonHandler(m.userAdminHandler.CreateBasicUser))
}
