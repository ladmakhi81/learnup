package course

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/course/handler"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	middleware    *middleware.Middleware
	courseHandler *handler.Handler
}

func NewModule(
	courseAdminHandler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		middleware:    middleware,
		courseHandler: courseAdminHandler,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	coursesApi := api.Group("/courses")

	coursesApi.Use(m.middleware.CheckAccessToken())

	coursesApi.POST("/", utils.JsonHandler(m.courseHandler.CreateCourse))
	coursesApi.GET("/page", utils.JsonHandler(m.courseHandler.GetCourses))
	coursesApi.GET("/:course-id/videos", utils.JsonHandler(m.courseHandler.GetVideosByCourseID))
	coursesApi.GET("/:course-id", utils.JsonHandler(m.courseHandler.GetCourseById))
	coursesApi.PATCH("/:course-id/verify", utils.JsonHandler(m.courseHandler.VerifyCourse))
	coursesApi.POST("/:course-id/like", utils.JsonHandler(m.courseHandler.Like))
	coursesApi.GET("/:course-id/likes", utils.JsonHandler(m.courseHandler.FetchLikes))
	coursesApi.POST("/:course-id/comment", utils.JsonHandler(m.courseHandler.CreateComment))
	coursesApi.DELETE("/comments/:comment-id", utils.JsonHandler(m.courseHandler.DeleteComment))
	coursesApi.POST("/:course-id/question", utils.JsonHandler(m.courseHandler.CreateQuestion))
	coursesApi.GET("/:course-id/questions", utils.JsonHandler(m.courseHandler.GetQuestions))
}
