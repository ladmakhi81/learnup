package user

import (
	"github.com/gin-gonic/gin"
	userApiHandler "github.com/ladmakhi81/learnup/internals/user/handler"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
	userHandler    *userApiHandler.Handler
}

func NewModule(
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
	userSvc userService.UserSvc,
	validationSvc contracts.Validation,
) *Module {
	return &Module{
		middleware:     middleware,
		translationSvc: translationSvc,
		userHandler: userApiHandler.NewHandler(
			userSvc,
			validationSvc,
			translationSvc,
		),
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	usersApi := api.Group("/users")
	usersApi.POST("/basic", utils.JsonHandler(m.translationSvc, m.userHandler.CreateBasicUser))
}
