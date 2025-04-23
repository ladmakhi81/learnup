package question

import (
	"github.com/gin-gonic/gin"
	questionHandler "github.com/ladmakhi81/learnup/internals/question/handler"
	questionService "github.com/ladmakhi81/learnup/internals/question/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	questionHandler *questionHandler.Handler
	middleware      *middleware.Middleware
	translationSvc  contracts.Translator
}

func NewModule(
	questionAnswerSvc questionService.QuestionAnswerService,
	validationSvc contracts.Validation,
	userSvc userService.UserSvc,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		questionHandler: questionHandler.NewHandler(questionAnswerSvc, translationSvc, validationSvc, userSvc),
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
