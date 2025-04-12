package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/notification/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	notificationAdminHandler *handler.Handler
	middlewares              *middleware.Middleware
}

func NewModule(
	notificationAdminHandler *handler.Handler,
	middlewares *middleware.Middleware,
) *Module {
	return &Module{
		notificationAdminHandler: notificationAdminHandler,
		middlewares:              middlewares,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	notificationsApi := api.Group("/notifications")

	notificationsApi.Use(m.middlewares.CheckAccessToken())

	notificationsApi.PATCH("/:notification-id/seen", utils.JsonHandler(m.notificationAdminHandler.SeenNotification))
	notificationsApi.GET("/page", utils.JsonHandler(m.notificationAdminHandler.GetNotificationsPage))
}
