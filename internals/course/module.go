package course

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/course/handler"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	middleware         *middleware.Middleware
	courseAdminHandler *handler.CourseAdminHandler
}

func NewModule(
	courseAdminHandler *handler.CourseAdminHandler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		middleware:         middleware,
		courseAdminHandler: courseAdminHandler,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	coursesApi := api.Group("/courses")
	adminCoursesApi := coursesApi.Group("/admin")

	adminCoursesApi.Use(m.middleware.CheckAccessToken())

	adminCoursesApi.POST("/", utils.JsonHandler(m.courseAdminHandler.CreateCourse))
}
