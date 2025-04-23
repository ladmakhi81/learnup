package comment

import (
	"github.com/gin-gonic/gin"
	commentHandler "github.com/ladmakhi81/learnup/internals/comment/handler"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	handler        *commentHandler.Handler
	middleware     *middleware.Middleware
	translationSvc contracts.Translator
}

func NewModule(
	commentSvc commentService.CommentService,
	validationSvc contracts.Validation,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		middleware:     middleware,
		translationSvc: translationSvc,
		handler:        commentHandler.NewHandler(commentSvc, translationSvc, validationSvc),
	}
}

func (m *Module) Register(api *gin.RouterGroup) {
	commentsApi := api.Group("/comments")
	commentsApi.Use(m.middleware.CheckAccessToken())
	commentsApi.GET("/page", utils.JsonHandler(m.translationSvc, m.handler.GetCommentsPageable))
}
