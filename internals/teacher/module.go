package teacher

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/middleware"
	courseHandler "github.com/ladmakhi81/learnup/internals/teacher/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	courseHandler *courseHandler.Handler
	middleware    *middleware.Middleware
}

func NewModule(
	courseHandler *courseHandler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		courseHandler: courseHandler,
		middleware:    middleware,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	teacherApi := api.Group("/teacher")

	teacherApi.Use(m.middleware.CheckAccessToken())

	teacherApi.POST("/course", utils.JsonHandler(m.courseHandler.CreateCourse))
	teacherApi.GET("/courses", utils.JsonHandler(m.courseHandler.FetchCourses))
}
