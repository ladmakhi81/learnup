package course

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/course/handler"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	middleware         *middleware.Middleware
	courseAdminHandler *handler.Handler
}

func NewModule(
	courseAdminHandler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		middleware:         middleware,
		courseAdminHandler: courseAdminHandler,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	coursesApi := api.Group("/courses")

	coursesApi.Use(m.middleware.CheckAccessToken())

	coursesApi.POST("/", utils.JsonHandler(m.courseAdminHandler.CreateCourse))
	coursesApi.GET("/page", utils.JsonHandler(m.courseAdminHandler.GetCourses))
	coursesApi.GET("/:course-id/videos", utils.JsonHandler(m.courseAdminHandler.GetVideosByCourseID))
	coursesApi.GET("/:course-id", utils.JsonHandler(m.courseAdminHandler.GetCourseById))
}
