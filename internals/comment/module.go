package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/comment/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	handler        *handler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	handler *handler.Handler,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		handler:        handler,
		middleware:     middleware,
		translationSvc: translationSvc,
	}
}

func (m *Module) Register(api *gin.RouterGroup) {
	commentsApi := api.Group("/comments")

	commentsApi.Use(m.middleware.CheckAccessToken())

	commentsApi.GET("/page", utils.JsonHandler(m.translationSvc, m.handler.GetCommentsPageable))
}
