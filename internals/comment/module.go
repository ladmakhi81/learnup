package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/comment/handler"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	handler    *handler.Handler
	middleware *middleware.Middleware
}

func NewModule(
	handler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		handler:    handler,
		middleware: middleware,
	}
}

func (m *Module) Register(api *gin.RouterGroup) {
	commentsApi := api.Group("/comments")

	commentsApi.Use(m.middleware.CheckAccessToken())
	
	commentsApi.GET("/page", utils.JsonHandler(m.handler.GetCommentsPageable))
}
