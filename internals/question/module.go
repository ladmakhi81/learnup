package question

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/question/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	questionHandler *handler.Handler
	middleware      *middleware.Middleware
}

func NewModule(
	questionHandler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		questionHandler: questionHandler,
		middleware:      middleware,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	questionsApi := api.Group("/questions")
	questionsApi.Use(m.middleware.CheckAccessToken())
	questionsApi.POST("/:question-id/answer", utils.JsonHandler(m.questionHandler.AnswerQuestion))
	questionsApi.GET("/:question-id/answers", utils.JsonHandler(m.questionHandler.GetQuestionAnswers))
}
