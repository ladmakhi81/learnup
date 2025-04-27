package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ladmakhi81/learnup/internals/auth"
	authService "github.com/ladmakhi81/learnup/internals/auth/service"
	"github.com/ladmakhi81/learnup/internals/cart"
	cartService "github.com/ladmakhi81/learnup/internals/cart/service"
	"github.com/ladmakhi81/learnup/internals/category"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	"github.com/ladmakhi81/learnup/internals/comment"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	"github.com/ladmakhi81/learnup/internals/course"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	forumService "github.com/ladmakhi81/learnup/internals/forum/service"
	likeService "github.com/ladmakhi81/learnup/internals/like/service"
	"github.com/ladmakhi81/learnup/internals/notification"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	"github.com/ladmakhi81/learnup/internals/order"
	orderService "github.com/ladmakhi81/learnup/internals/order/service"
	"github.com/ladmakhi81/learnup/internals/payment"
	paymentService "github.com/ladmakhi81/learnup/internals/payment/service"
	"github.com/ladmakhi81/learnup/internals/question"
	questionService "github.com/ladmakhi81/learnup/internals/question/service"
	"github.com/ladmakhi81/learnup/internals/teacher"
	teacherService "github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/internals/transaction"
	transactionService "github.com/ladmakhi81/learnup/internals/transaction/service"
	"github.com/ladmakhi81/learnup/internals/tus"
	tusHookService "github.com/ladmakhi81/learnup/internals/tus/service"
	"github.com/ladmakhi81/learnup/internals/user"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/internals/video"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/internals/video/workflow"
	"github.com/ladmakhi81/learnup/internals/websocket"
	"github.com/ladmakhi81/learnup/pkg/ffmpeg/v1"
	"github.com/ladmakhi81/learnup/pkg/i18n/v2"
	"github.com/ladmakhi81/learnup/pkg/jwt/v5"
	"github.com/ladmakhi81/learnup/pkg/koanf"
	"github.com/ladmakhi81/learnup/pkg/logrus/v1"
	"github.com/ladmakhi81/learnup/pkg/minio/v7"
	"github.com/ladmakhi81/learnup/pkg/redis/v6"
	restyv2 "github.com/ladmakhi81/learnup/pkg/resty/v2"
	stripev82 "github.com/ladmakhi81/learnup/pkg/stripe/v82"
	"github.com/ladmakhi81/learnup/pkg/temporal"
	"github.com/ladmakhi81/learnup/pkg/temporal/v1"
	"github.com/ladmakhi81/learnup/pkg/validator/v10"
	zarinpalv1 "github.com/ladmakhi81/learnup/pkg/zarinpal/v1"
	zibalv1 "github.com/ladmakhi81/learnup/pkg/zibal/v1"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"log"

	_ "github.com/ladmakhi81/learnup/docs"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Learnup
// @version         1.0
// @BasePath  /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// config file loader
	koanfConfigProvider := koanf.NewKoanfEnvSvc()
	config, configErr := koanfConfigProvider.LoadLearnUp()
	if configErr != nil {
		log.Fatalf("load learn up config failed: %v", configErr)
	}

	// temporal
	temporalSvc := temporalv1.NewTemporalSvc(config)
	if err := temporalSvc.Init(); err != nil {
		log.Fatalf("temporal throw error: %v", err)
	}

	// database
	dbClient := db.NewDatabase(config)
	if err := dbClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// http handler
	server := gin.Default()
	port := fmt.Sprintf(":%d", config.App.Port)
	api := server.Group("/api")

	// swagger documentation
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// websocket
	wsManager := websocket.NewWsManager()

	unitOfWork := db.NewUnitOfWork(dbClient.Core)

	// svcs
	logrusSvc := logrusv1.NewLogrusLoggerSvc()
	minioSvc, minioSvcErr := miniov7.NewMinioClientSvc(config)
	if minioSvcErr != nil {
		log.Fatalln(minioSvcErr)
	}
	i18nTranslatorSvc, i18nErr := i18nv2.NewI18nTranslatorSvc()
	if i18nErr != nil {
		log.Fatalln(i18nErr)
	}
	redisSvc := redisv6.NewRedisClientSvc(config)
	tokenSvc := jwtv5.NewJwtSvc(config, redisSvc)
	userSvc := userService.NewUserSvc(unitOfWork)
	notificationSvc := notificationService.NewNotificationSvc(unitOfWork)
	validationSvc := validatorv10.NewValidatorSvc(validator.New(), i18nTranslatorSvc)
	authSvc := authService.NewAuthSvc(redisSvc, tokenSvc, unitOfWork)
	categorySvc := categoryService.NewCategorySvc(unitOfWork)
	courseSvc := courseService.NewCourseSvc(unitOfWork)
	forumSvc := forumService.NewForumService(unitOfWork)
	ffmpegSvc := ffmpegv1.NewFfmpegSvc()
	videoSvc := videoService.NewVideoSvc(unitOfWork, minioSvc, ffmpegSvc, logrusSvc)
	videoWorkflowSvc := workflow.NewVideoWorkflowImpl(videoSvc, temporalSvc, courseSvc)
	tusHookSvc := tusHookService.NewTusServiceImpl(videoSvc, logrusSvc, temporalSvc, videoWorkflowSvc)
	teacherCourseSvc := teacherService.NewTeacherCourseService(unitOfWork)
	teacherVideoSvc := teacherService.NewTeacherVideoSvc(unitOfWork)
	teacherCommentSvc := teacherService.NewTeacherCommentSvc(unitOfWork)
	commentSvc := commentService.NewCommentSvc(unitOfWork)
	likeSvc := likeService.NewLikeSvc(unitOfWork)
	cartSvc := cartService.NewCartSvc(unitOfWork)
	questionSvc := questionService.NewQuestionSvc(unitOfWork)
	questionAnswerSvc := questionService.NewQuestionAnswerSvc(unitOfWork)
	teacherQuestionSvc := teacherService.NewTeacherQuestionSvc(unitOfWork)
	restyHttpClient := restyv2.NewRestyHttpSvc()
	zarinpalSvc := zarinpalv1.NewZarinpalClient(restyHttpClient, config)
	zibalSvc := zibalv1.NewZibalClient(restyHttpClient, config)
	stripeSvc, stripeSvcErr := stripev82.NewStripeClient(config)
	if stripeSvcErr != nil {
		panic("stripe client error occured")
	}
	transactionSvc := transactionService.NewTransactionSvc(unitOfWork)
	paymentSvc := paymentService.NewPaymentService(unitOfWork, zarinpalSvc, zibalSvc, stripeSvc, config)
	orderSvc := orderService.NewOrderService(unitOfWork, paymentSvc)

	// middlewares
	middlewares := middleware.NewMiddleware(tokenSvc)

	// modules
	userModule := user.NewModule(middlewares, i18nTranslatorSvc, userSvc, validationSvc)
	authModule := auth.NewModule(authSvc, validationSvc, i18nTranslatorSvc)
	categoryModule := category.NewModule(categorySvc, middlewares, i18nTranslatorSvc, validationSvc)
	courseModule := course.NewModule(courseSvc, validationSvc, videoSvc, likeSvc, commentSvc, questionSvc, userSvc, forumSvc, middlewares, i18nTranslatorSvc)
	tusModule := tus.NewModule(tusHookSvc, i18nTranslatorSvc)
	videoModule := video.NewModule(userSvc, videoSvc, validationSvc, middlewares, i18nTranslatorSvc)
	notificationModule := notification.NewModule(notificationSvc, middlewares, i18nTranslatorSvc)
	teacherModule := teacher.NewModule(teacherCourseSvc, teacherVideoSvc, teacherCommentSvc, teacherQuestionSvc, validationSvc, userSvc, middlewares, i18nTranslatorSvc)
	commentModule := comment.NewModule(commentSvc, validationSvc, middlewares, i18nTranslatorSvc)
	questionModule := question.NewModule(questionAnswerSvc, validationSvc, userSvc, middlewares, i18nTranslatorSvc)
	cartModule := cart.NewModule(userSvc, cartSvc, validationSvc, middlewares, i18nTranslatorSvc)
	orderModule := order.NewModule(orderSvc, validationSvc, userSvc, middlewares, i18nTranslatorSvc)
	paymentModule := payment.NewModule(paymentSvc, middlewares, i18nTranslatorSvc)
	transactionModule := transaction.NewModule(transactionSvc, middlewares, i18nTranslatorSvc)

	server.GET("/ws", websocket.Handler(wsManager, tokenSvc))

	// workers
	if err := temporalSvc.AddWorker(
		temporal.ADD_NEW_COURSE_VIDEO_QUEUE,
		videoWorkflowSvc.AddNewCourseVideoWorkflow,
		videoSvc.CalculateDuration,
		videoSvc.Encode,
		videoSvc.UpdateURLAndDuration,
		videoSvc.CreateCompleteUploadVideoNotification,
	); err != nil {
		log.Printf("Error in add worker: %+v", err)
	}

	if err := temporalSvc.AddWorker(
		temporal.SET_INTRODUCTION_COURSE_QUEUE,
		videoWorkflowSvc.AddIntroductionVideoWorkflow,
		videoSvc.Encode,
		courseSvc.UpdateIntroductionURL,
		courseSvc.CreateCompleteIntroductionVideoNotification,
	); err != nil {
		log.Printf("Error in add worker: %+v", err)
	}

	// register module
	userModule.Register(api)
	authModule.Register(api)
	categoryModule.Register(api)
	courseModule.Register(api)
	tusModule.Register(api)
	videoModule.Register(api)
	notificationModule.Register(api)
	teacherModule.Register(api)
	commentModule.Register(api)
	questionModule.Register(api)
	cartModule.Register(api)
	orderModule.Register(api)
	paymentModule.Register(api)
	transactionModule.Register(api)

	log.Printf("the server running on %s \n", port)

	// run http handler
	log.Fatalln(server.Run(port))
}
