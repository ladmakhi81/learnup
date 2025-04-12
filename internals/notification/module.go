package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/notification/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	notificationAdminHandler *handler.NotificationAdminHandler
	middlewares              *middleware.Middleware
}

func NewModule(
	notificationAdminHandler *handler.NotificationAdminHandler,
	middlewares *middleware.Middleware,
) *Module {
	return &Module{
		notificationAdminHandler: notificationAdminHandler,
		middlewares:              middlewares,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	notificationsApi := api.Group("/notifications")
	adminNotificationsApi := notificationsApi.Group("/admin")

	adminNotificationsApi.Use(m.middlewares.CheckAccessToken())

	adminNotificationsApi.PATCH("/:notification-id/seen", utils.JsonHandler(m.notificationAdminHandler.SeenNotification))
	adminNotificationsApi.GET("/page", utils.JsonHandler(m.notificationAdminHandler.GetNotificationsPage))
}
