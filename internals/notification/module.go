package notification

import (
	"github.com/gin-gonic/gin"
	notificationHandler "github.com/ladmakhi81/learnup/internals/notification/handler"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	notificationAdminHandler *notificationHandler.Handler
	middlewares              *middleware.Middleware
	translationSvc           contracts.Translator
}

func NewModule(
	notificationSvc notificationService.NotificationService,
	middlewares *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		notificationAdminHandler: notificationHandler.NewHandler(
			notificationSvc,
			translationSvc,
		),
		middlewares:    middlewares,
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	notificationsApi := api.Group("/notifications")
	notificationsApi.Use(m.middlewares.CheckAccessToken())
	notificationsApi.PATCH("/:notification-id/seen", utils.JsonHandler(m.translationSvc, m.notificationAdminHandler.SeenNotification))
	notificationsApi.GET("/page", utils.JsonHandler(m.translationSvc, m.notificationAdminHandler.GetNotificationsPage))
}
