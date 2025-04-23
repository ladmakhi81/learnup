package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/notification/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	notificationAdminHandler *handler.Handler
	middlewares              *middleware.Middleware
	translationSvc           contracts.Translator
}

func NewModule(
	notificationAdminHandler *handler.Handler,
	middlewares *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		notificationAdminHandler: notificationAdminHandler,
		middlewares:              middlewares,
		translationSvc:           translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	notificationsApi := api.Group("/notifications")

	notificationsApi.Use(m.middlewares.CheckAccessToken())

	notificationsApi.PATCH("/:notification-id/seen", utils.JsonHandler(m.translationSvc, m.notificationAdminHandler.SeenNotification))
	notificationsApi.GET("/page", utils.JsonHandler(m.translationSvc, m.notificationAdminHandler.GetNotificationsPage))
}
