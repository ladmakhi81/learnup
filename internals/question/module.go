package question

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/question/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	questionHandler *handler.Handler
	middleware      *middleware.Middleware
	translationSvc  contracts.Translator
}

func NewModule(
	questionHandler *handler.Handler,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		questionHandler: questionHandler,
		middleware:      middleware,
		translationSvc:  translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	questionsApi := api.Group("/questions")
	questionsApi.Use(m.middleware.CheckAccessToken())
	questionsApi.POST("/:question-id/answer", utils.JsonHandler(m.translationSvc, m.questionHandler.AnswerQuestion))
	questionsApi.GET("/:question-id/answers", utils.JsonHandler(m.translationSvc, m.questionHandler.GetQuestionAnswers))
}
