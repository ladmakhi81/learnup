package teacher

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/teacher/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	courseHandler   *handler.CourseHandler
	middleware      *middleware.Middleware
	videoHandler    *handler.VideoHandler
	commentHandler  *handler.CommentHandler
	questionHandler *handler.QuestionHandler
	translationSvc  contracts.Translator
}

func NewModule(
	courseHandler *handler.CourseHandler,
	videoHandler *handler.VideoHandler,
	middleware *middleware.Middleware,
	commentHandler *handler.CommentHandler,
	questionHandler *handler.QuestionHandler,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		courseHandler:   courseHandler,
		middleware:      middleware,
		videoHandler:    videoHandler,
		commentHandler:  commentHandler,
		questionHandler: questionHandler,
		translationSvc:  translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	teacherApi := api.Group("/teacher")

	teacherApi.Use(m.middleware.CheckAccessToken())

	teacherApi.POST("/course", utils.JsonHandler(m.translationSvc, m.courseHandler.CreateCourse))
	teacherApi.GET("/courses", utils.JsonHandler(m.translationSvc, m.courseHandler.FetchCourses))
	teacherApi.POST("/video", utils.JsonHandler(m.translationSvc, m.videoHandler.AddVideoToCourse))
	teacherApi.GET("/comments/:course-id", utils.JsonHandler(m.translationSvc, m.commentHandler.GetPageableCommentByCourseId))
	teacherApi.GET("/questions", utils.JsonHandler(m.translationSvc, m.questionHandler.GetQuestions))
}
