package course

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/course/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	middleware     *middleware.Middleware
	courseHandler  *handler.Handler
	translationSvc contracts.Translator
}

func NewModule(
	courseAdminHandler *handler.Handler,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		middleware:     middleware,
		courseHandler:  courseAdminHandler,
		translationSvc: translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	coursesApi := api.Group("/courses")

	coursesApi.Use(m.middleware.CheckAccessToken())

	coursesApi.POST("/", utils.JsonHandler(m.translationSvc, m.courseHandler.CreateCourse))
	coursesApi.GET("/page", utils.JsonHandler(m.translationSvc, m.courseHandler.GetCourses))
	coursesApi.GET("/:course-id/videos", utils.JsonHandler(m.translationSvc, m.courseHandler.GetVideosByCourseID))
	coursesApi.GET("/:course-id", utils.JsonHandler(m.translationSvc, m.courseHandler.GetCourseById))
	coursesApi.PATCH("/:course-id/verify", utils.JsonHandler(m.translationSvc, m.courseHandler.VerifyCourse))
	coursesApi.POST("/:course-id/like", utils.JsonHandler(m.translationSvc, m.courseHandler.Like))
	coursesApi.GET("/:course-id/likes", utils.JsonHandler(m.translationSvc, m.courseHandler.FetchLikes))
	coursesApi.POST("/:course-id/comment", utils.JsonHandler(m.translationSvc, m.courseHandler.CreateComment))
	coursesApi.DELETE("/comments/:comment-id", utils.JsonHandler(m.translationSvc, m.courseHandler.DeleteComment))
	coursesApi.POST("/:course-id/question", utils.JsonHandler(m.translationSvc, m.courseHandler.CreateQuestion))
	coursesApi.GET("/:course-id/questions", utils.JsonHandler(m.translationSvc, m.courseHandler.GetQuestions))
}
