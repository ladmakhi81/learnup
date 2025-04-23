package teacher

import (
	"github.com/gin-gonic/gin"
	teacherHandler "github.com/ladmakhi81/learnup/internals/teacher/handler"
	teacherService "github.com/ladmakhi81/learnup/internals/teacher/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	courseHandler   *teacherHandler.CourseHandler
	middleware      *middleware.Middleware
	videoHandler    *teacherHandler.VideoHandler
	commentHandler  *teacherHandler.CommentHandler
	questionHandler *teacherHandler.QuestionHandler
	translationSvc  contracts.Translator
}

func NewModule(
	teacherCourseSvc teacherService.TeacherCourseService,
	teacherVideoSvc teacherService.TeacherVideoService,
	teacherCommentSvc teacherService.TeacherCommentService,
	teacherQuestionSvc teacherService.TeacherQuestionService,
	validationSvc contracts.Validation,
	userSvc userService.UserSvc,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		courseHandler: teacherHandler.NewCourseHandler(
			teacherCourseSvc,
			validationSvc,
			translationSvc,
			userSvc,
		),
		middleware: middleware,
		videoHandler: teacherHandler.NewVideoHandler(
			teacherVideoSvc,
			translationSvc,
			validationSvc,
		),
		commentHandler: teacherHandler.NewCommentHandler(
			teacherCommentSvc,
			translationSvc,
			userSvc,
		),
		questionHandler: teacherHandler.NewQuestionHandler(
			translationSvc,
			teacherQuestionSvc,
			userSvc,
		),
		translationSvc: translationSvc,
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
