package teacher

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/teacher/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	courseHandler *handler.CourseHandler
	middleware    *middleware.Middleware
	videoHandler  *handler.VideoHandler
}

func NewModule(
	courseHandler *handler.CourseHandler,
	videoHandler *handler.VideoHandler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		courseHandler: courseHandler,
		middleware:    middleware,
		videoHandler:  videoHandler,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	teacherApi := api.Group("/teacher")

	teacherApi.Use(m.middleware.CheckAccessToken())

	teacherApi.POST("/course", utils.JsonHandler(m.courseHandler.CreateCourse))
	teacherApi.GET("/courses", utils.JsonHandler(m.courseHandler.FetchCourses))
	teacherApi.POST("/video", utils.JsonHandler(m.videoHandler.AddVideoToCourse))
}
