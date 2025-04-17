package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/auth"
	authApiHandler "github.com/ladmakhi81/learnup/internals/auth/handler"
	authService "github.com/ladmakhi81/learnup/internals/auth/service"
	"github.com/ladmakhi81/learnup/internals/category"
	categoryApiHandler "github.com/ladmakhi81/learnup/internals/category/handler"
	categoryRepository "github.com/ladmakhi81/learnup/internals/category/repo"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	"github.com/ladmakhi81/learnup/internals/comment"
	commentApiHandler "github.com/ladmakhi81/learnup/internals/comment/handler"
	commentRepository "github.com/ladmakhi81/learnup/internals/comment/repo"
	commentService "github.com/ladmakhi81/learnup/internals/comment/service"
	"github.com/ladmakhi81/learnup/internals/course"
	courseApiHandler "github.com/ladmakhi81/learnup/internals/course/handler"
	courseRepository "github.com/ladmakhi81/learnup/internals/course/repo"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	likeRepository "github.com/ladmakhi81/learnup/internals/like/repo"
	likeService "github.com/ladmakhi81/learnup/internals/like/service"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/notification"
	notificationApiHandler "github.com/ladmakhi81/learnup/internals/notification/handler"
	notificationRepository "github.com/ladmakhi81/learnup/internals/notification/repo"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	"github.com/ladmakhi81/learnup/internals/teacher"
	teacherApiHandler "github.com/ladmakhi81/learnup/internals/teacher/handler"
	teacherService "github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/internals/tus"
	tusHookApiHandler "github.com/ladmakhi81/learnup/internals/tus/handler"
	tusHookService "github.com/ladmakhi81/learnup/internals/tus/service"
	"github.com/ladmakhi81/learnup/internals/user"
	userApiHandler "github.com/ladmakhi81/learnup/internals/user/handler"
	userRepository "github.com/ladmakhi81/learnup/internals/user/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/internals/video"
	videoApiHandler "github.com/ladmakhi81/learnup/internals/video/handler"
	videoRepository "github.com/ladmakhi81/learnup/internals/video/repo"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/internals/video/workflow"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/pkg/ffmpeg/v1"
	"github.com/ladmakhi81/learnup/pkg/i18n/v2"
	"github.com/ladmakhi81/learnup/pkg/jwt/v5"
	"github.com/ladmakhi81/learnup/pkg/koanf"
	"github.com/ladmakhi81/learnup/pkg/logrus/v1"
	"github.com/ladmakhi81/learnup/pkg/minio/v7"
	"github.com/ladmakhi81/learnup/pkg/redis/v6"
	"github.com/ladmakhi81/learnup/pkg/temporal"
	"github.com/ladmakhi81/learnup/pkg/temporal/v1"
	"github.com/ladmakhi81/learnup/pkg/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"path"

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

	// minio
	minioClient, minioClientErr := SetupMinio(config)
	if minioClientErr != nil {
		log.Fatalf("Failed to connect minio: %v", minioClientErr)
	}

	// redis
	redisClient := SetupRedis(config)

	// database
	dbClient := db.NewDatabase(config)
	if err := dbClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// http handler
	server := gin.Default()
	port := fmt.Sprintf(":%d", config.App.Port)
	api := server.Group("/api")

	localizer, localizerErr := SetupLocalizer()
	if localizerErr != nil {
		log.Fatalf("Failed to initialize localizer: %v\n", localizerErr)
	}

	// swagger documentation
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// repos
	categoryRepo := categoryRepository.NewCategoryRepoImpl(dbClient)
	userRepo := userRepository.NewUserRepoImpl(dbClient)
	courseRepo := courseRepository.NewCourseRepoImpl(dbClient)
	videoRepo := videoRepository.NewVideoRepoImpl(dbClient)
	notificationRepo := notificationRepository.NewNotificationRepoImpl(dbClient)
	commentRepo := commentRepository.NewCommentRepoImpl(dbClient)
	likeRepo := likeRepository.NewLikeRepoImpl(dbClient)

	// svcs
	logrusSvc := logrusv1.NewLogrusLoggerSvc()
	minioSvc := miniov7.NewMinioClientSvc(minioClient, config.Minio.Region)
	i18nTranslatorSvc := i18nv2.NewI18nTranslatorSvc(localizer)
	redisSvc := redisv6.NewRedisClientSvc(redisClient)
	tokenSvc := jwtv5.NewJwtSvc(config)
	userSvc := userService.NewUserSvcImpl(userRepo, i18nTranslatorSvc)
	notificationSvc := notificationService.NewNotificationServiceImpl(notificationRepo, userSvc, i18nTranslatorSvc)
	validationSvc := validatorv10.NewValidatorSvc(validator.New(), i18nTranslatorSvc)
	authSvc := authService.NewAuthServiceImpl(userSvc, redisSvc, tokenSvc, i18nTranslatorSvc)
	categorySvc := categoryService.NewCategoryServiceImpl(categoryRepo, i18nTranslatorSvc)
	courseSvc := courseService.NewCourseServiceImpl(courseRepo, i18nTranslatorSvc, userSvc, categorySvc, notificationSvc)
	ffmpegSvc := ffmpegv1.NewFfmpegSvc()
	videoSvc := videoService.NewVideoServiceImpl(minioSvc, ffmpegSvc, logrusSvc, courseSvc, videoRepo, notificationSvc, i18nTranslatorSvc, userSvc)
	videoWorkflowSvc := workflow.NewVideoWorkflowImpl(videoSvc, temporalSvc, courseSvc)
	tusHookSvc := tusHookService.NewTusServiceImpl(videoSvc, logrusSvc, temporalSvc, videoWorkflowSvc)
	teacherCourseSvc := teacherService.NewTeacherCourseServiceImpl(courseSvc, categorySvc, userSvc, courseRepo, i18nTranslatorSvc)
	teacherVideoSvc := teacherService.NewTeacherVideoServiceImpl(videoSvc, i18nTranslatorSvc, courseSvc, videoRepo)
	teacherCommentSvc := teacherService.NewTeacherCommentServiceImpl(userSvc, courseSvc, commentRepo, i18nTranslatorSvc)
	commentSvc := commentService.NewCommentServiceImpl(commentRepo, userSvc, courseSvc, i18nTranslatorSvc)
	likeSvc := likeService.NewLikeServiceImpl(likeRepo, userSvc, i18nTranslatorSvc, courseSvc)

	// middlewares
	middlewares := middleware.NewMiddleware(tokenSvc, redisSvc)

	// handlers
	userHandler := userApiHandler.NewHandler(userSvc, validationSvc, i18nTranslatorSvc)
	authHandler := authApiHandler.NewHandler(authSvc, validationSvc, i18nTranslatorSvc)
	categoryHandler := categoryApiHandler.NewHandler(categorySvc, i18nTranslatorSvc, validationSvc)
	courseHandler := courseApiHandler.NewHandler(courseSvc, validationSvc, i18nTranslatorSvc, videoSvc, likeSvc, commentSvc)
	videoHandler := videoApiHandler.NewHandler(validationSvc, videoSvc, i18nTranslatorSvc)
	tusHandler := tusHookApiHandler.NewTusHookHandler(tusHookSvc)
	notificationHandler := notificationApiHandler.NewHandler(notificationSvc, i18nTranslatorSvc)
	teacherCourseHandler := teacherApiHandler.NewCourseHandler(teacherCourseSvc, validationSvc, i18nTranslatorSvc)
	teacherVideoHandler := teacherApiHandler.NewVideoHandler(teacherVideoSvc, i18nTranslatorSvc, validationSvc)
	teacherCommentHandler := teacherApiHandler.NewCommentHandler(teacherCommentSvc, i18nTranslatorSvc)
	commentHandler := commentApiHandler.NewHandler(commentSvc, i18nTranslatorSvc, validationSvc)

	// modules
	userModule := user.NewModule(userHandler, middlewares)
	authModule := auth.NewModule(authHandler)
	categoryModule := category.NewModule(categoryHandler, middlewares)
	courseModule := course.NewModule(courseHandler, middlewares)
	tusModule := tus.NewModule(tusHandler)
	videoModule := video.NewModule(videoHandler, middlewares)
	notificationModule := notification.NewModule(notificationHandler, middlewares)
	teacherModule := teacher.NewModule(teacherCourseHandler, teacherVideoHandler, middlewares, teacherCommentHandler)
	commentModule := comment.NewModule(commentHandler, middlewares)

	// workers
	if err := temporalSvc.AddWorker(
		temporal.ADD_NEW_COURSE_VIDEO_QUEUE,
		videoWorkflowSvc.AddNewCourseVideoWorkflow,
		videoSvc.CalculateDuration,
		videoSvc.Encode,
		videoSvc.UpdateURLAndDuration,
		notificationSvc.Create,
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

	log.Printf("the server running on %s \n", port)

	// run http handler
	log.Fatalln(server.Run(port))
}

func SetupMinio(config *dtos.EnvConfig) (*minio.Client, error) {
	endpoint := config.Minio.URL
	accessKey := config.Minio.AccessKey
	secretKey := config.Minio.SecretKey
	region := config.Minio.Region
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
		Region: region,
	})
}

func SetupRedis(config *dtos.EnvConfig) *redis.Client {
	host := config.Redis.Host
	port := config.Redis.Port
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})
}

func SetupLocalizer() (*i18n.Localizer, error) {
	locales := map[string]struct {
		langTag  language.Tag
		langText string
	}{
		"fa": {
			langTag:  language.Persian,
			langText: "fa",
		},
		"en": {
			langTag:  language.English,
			langText: "en",
		},
	}
	defaultLocale := "fa"
	localizationBundle := i18n.NewBundle(locales[defaultLocale].langTag)
	localizationBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	rootDir, rootDirErr := os.Getwd()
	if rootDirErr != nil {
		return nil, rootDirErr
	}
	translationFolderPath := path.Join(rootDir, "translations")
	localizationBundle.MustLoadMessageFile(path.Join(translationFolderPath, "fa.json"))
	localizationBundle.MustLoadMessageFile(path.Join(translationFolderPath, "en.json"))
	return i18n.NewLocalizer(localizationBundle, locales[defaultLocale].langText), nil
}
