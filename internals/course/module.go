package course

import (
	"github.com/gin-gonic/gin"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	courseHandler "github.com/ladmakhi81/learnup/internals/course/handler"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	forumService "github.com/ladmakhi81/learnup/internals/forum/service"
	likeService "github.com/ladmakhi81/learnup/internals/like/service"
	questionService "github.com/ladmakhi81/learnup/internals/question/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	middleware     *middleware.Middleware
	courseHandler  *courseHandler.Handler
	translationSvc contracts.Translator
}

func NewModule(
	courseSvc courseService.CourseService,
	validationSvc contracts.Validation,
	videoSvc videoService.VideoService,
	likeSvc likeService.LikeService,
	commentSvc commentService.CommentService,
	questionSvc questionService.QuestionService,
	userSvc userService.UserSvc,
	forumSvc forumService.ForumService,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		middleware:     middleware,
		translationSvc: translationSvc,
		courseHandler: courseHandler.NewHandler(
			courseSvc,
			validationSvc,
			translationSvc,
			videoSvc,
			likeSvc,
			commentSvc,
			questionSvc,
			userSvc,
			forumSvc,
		),
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
	coursesApi.GET("/:course-id/forum", utils.JsonHandler(m.translationSvc, m.courseHandler.GetForumByCourseID))
}
