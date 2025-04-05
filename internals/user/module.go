package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/user/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	userAdminHandler *handler.UserAdminHandler
}

func NewModule(
	userAdminHandler *handler.UserAdminHandler,
) *Module {
	return &Module{
		userAdminHandler: userAdminHandler,
	}
}

func (m *Module) Register(api *gin.RouterGroup) {
	usersApi := api.Group("/users")
	usersAdminApi := usersApi.Group("/admin")

	usersAdminApi.POST("/", utils.JsonHandler(m.userAdminHandler.CreateUser))
}
