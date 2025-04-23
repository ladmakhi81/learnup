package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ladmakhi81/learnup/internals/auth"
	authApiHandler "github.com/ladmakhi81/learnup/internals/auth/handler"
	authService "github.com/ladmakhi81/learnup/internals/auth/service"
	"github.com/ladmakhi81/learnup/internals/cart"
	cartApiHandler "github.com/ladmakhi81/learnup/internals/cart/handler"
	cartService "github.com/ladmakhi81/learnup/internals/cart/service"
	"github.com/ladmakhi81/learnup/internals/category"
	categoryApiHandler "github.com/ladmakhi81/learnup/internals/category/handler"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	"github.com/ladmakhi81/learnup/internals/comment"
	commentApiHandler "github.com/ladmakhi81/learnup/internals/comment/handler"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	"github.com/ladmakhi81/learnup/internals/course"
	courseApiHandler "github.com/ladmakhi81/learnup/internals/course/handler"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	likeService "github.com/ladmakhi81/learnup/internals/like/service"
	"github.com/ladmakhi81/learnup/internals/notification"
	notificationApiHandler "github.com/ladmakhi81/learnup/internals/notification/handler"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	"github.com/ladmakhi81/learnup/internals/order"
	orderApiHandler "github.com/ladmakhi81/learnup/internals/order/handler"
	orderService "github.com/ladmakhi81/learnup/internals/order/service"
	"github.com/ladmakhi81/learnup/internals/payment"
	paymentApiHandler "github.com/ladmakhi81/learnup/internals/payment/handler"
	paymentService "github.com/ladmakhi81/learnup/internals/payment/service"
	"github.com/ladmakhi81/learnup/internals/question"
	questionApiHandler "github.com/ladmakhi81/learnup/internals/question/handler"
	questionService "github.com/ladmakhi81/learnup/internals/question/service"
	"github.com/ladmakhi81/learnup/internals/teacher"
	teacherApiHandler "github.com/ladmakhi81/learnup/internals/teacher/handler"
	teacherService "github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/internals/transaction"
	transactionApiHandler "github.com/ladmakhi81/learnup/internals/transaction/handler"
	transactionService "github.com/ladmakhi81/learnup/internals/transaction/service"
	"github.com/ladmakhi81/learnup/internals/tus"
	tusHookApiHandler "github.com/ladmakhi81/learnup/internals/tus/handler"
	tusHookService "github.com/ladmakhi81/learnup/internals/tus/service"
	"github.com/ladmakhi81/learnup/internals/user"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/internals/video"
	videoApiHandler "github.com/ladmakhi81/learnup/internals/video/handler"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/internals/video/workflow"
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
	tokenSvc := jwtv5.NewJwtSvc(config)
	userSvc := userService.NewUserSvc(unitOfWork)
	notificationSvc := notificationService.NewNotificationSvc(unitOfWork)
	validationSvc := validatorv10.NewValidatorSvc(validator.New(), i18nTranslatorSvc)
	authSvc := authService.NewAuthSvc(redisSvc, tokenSvc, unitOfWork)
	categorySvc := categoryService.NewCategorySvc(unitOfWork)
	courseSvc := courseService.NewCourseSvc(unitOfWork)
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
	middlewares := middleware.NewMiddleware(tokenSvc, redisSvc)

	// handlers
	authHandler := authApiHandler.NewHandler(authSvc, validationSvc, i18nTranslatorSvc)
	categoryHandler := categoryApiHandler.NewHandler(categorySvc, i18nTranslatorSvc, validationSvc)
	courseHandler := courseApiHandler.NewHandler(courseSvc, validationSvc, i18nTranslatorSvc, videoSvc, likeSvc, commentSvc, questionSvc, userSvc)
	videoHandler := videoApiHandler.NewHandler(validationSvc, videoSvc, i18nTranslatorSvc, userSvc)
	tusHandler := tusHookApiHandler.NewTusHookHandler(tusHookSvc)
	notificationHandler := notificationApiHandler.NewHandler(notificationSvc, i18nTranslatorSvc)
	teacherCourseHandler := teacherApiHandler.NewCourseHandler(teacherCourseSvc, validationSvc, i18nTranslatorSvc, userSvc)
	teacherVideoHandler := teacherApiHandler.NewVideoHandler(teacherVideoSvc, i18nTranslatorSvc, validationSvc)
	teacherCommentHandler := teacherApiHandler.NewCommentHandler(teacherCommentSvc, i18nTranslatorSvc, userSvc)
	teacherQuestionHandler := teacherApiHandler.NewQuestionHandler(i18nTranslatorSvc, teacherQuestionSvc, userSvc)
	commentHandler := commentApiHandler.NewHandler(commentSvc, i18nTranslatorSvc, validationSvc)
	questionHandler := questionApiHandler.NewHandler(questionAnswerSvc, i18nTranslatorSvc, validationSvc, userSvc)
	cartHandler := cartApiHandler.NewHandler(i18nTranslatorSvc, validationSvc, cartSvc, userSvc)
	orderHandler := orderApiHandler.NewHandler(orderSvc, i18nTranslatorSvc, validationSvc, userSvc)
	paymentHandler := paymentApiHandler.NewHandler(paymentSvc)
	transactionHandler := transactionApiHandler.NewHandler(transactionSvc)

	// modules
	userModule := user.NewModule(middlewares, i18nTranslatorSvc, userSvc, validationSvc)
	authModule := auth.NewModule(i18nTranslatorSvc, authHandler)
	categoryModule := category.NewModule(categoryHandler, middlewares, i18nTranslatorSvc)
	courseModule := course.NewModule(courseHandler, middlewares, i18nTranslatorSvc)
	tusModule := tus.NewModule(tusHandler, i18nTranslatorSvc)
	videoModule := video.NewModule(videoHandler, middlewares, i18nTranslatorSvc)
	notificationModule := notification.NewModule(notificationHandler, middlewares, i18nTranslatorSvc)
	teacherModule := teacher.NewModule(teacherCourseHandler, teacherVideoHandler, middlewares, teacherCommentHandler, teacherQuestionHandler, i18nTranslatorSvc)
	commentModule := comment.NewModule(commentHandler, middlewares, i18nTranslatorSvc)
	questionModule := question.NewModule(questionHandler, middlewares, i18nTranslatorSvc)
	cartModule := cart.NewModule(cartHandler, middlewares, i18nTranslatorSvc)
	orderModule := order.NewModule(orderHandler, middlewares, i18nTranslatorSvc)
	paymentModule := payment.NewModule(paymentHandler, middlewares, i18nTranslatorSvc)
	transactionModule := transaction.NewModule(transactionHandler, middlewares, i18nTranslatorSvc)

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
