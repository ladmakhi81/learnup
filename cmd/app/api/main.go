package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/auth"
	authHandler "github.com/ladmakhi81/learnup/internals/auth/handler"
	authService "github.com/ladmakhi81/learnup/internals/auth/service"
	"github.com/ladmakhi81/learnup/internals/category"
	categoryHandler "github.com/ladmakhi81/learnup/internals/category/handler"
	categoryRepository "github.com/ladmakhi81/learnup/internals/category/repo"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	"github.com/ladmakhi81/learnup/internals/course"
	courseHandler "github.com/ladmakhi81/learnup/internals/course/handler"
	courseRepository "github.com/ladmakhi81/learnup/internals/course/repo"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/internals/notification"
	notificationHandler "github.com/ladmakhi81/learnup/internals/notification/handler"
	notificationRepository "github.com/ladmakhi81/learnup/internals/notification/repo"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	"github.com/ladmakhi81/learnup/internals/tus"
	tusHookHandler "github.com/ladmakhi81/learnup/internals/tus/handler"
	tusHookService "github.com/ladmakhi81/learnup/internals/tus/service"
	"github.com/ladmakhi81/learnup/internals/user"
	userHandler "github.com/ladmakhi81/learnup/internals/user/handler"
	userRepository "github.com/ladmakhi81/learnup/internals/user/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/internals/video"
	videoHandler "github.com/ladmakhi81/learnup/internals/video/handler"
	videoRepository "github.com/ladmakhi81/learnup/internals/video/repo"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	redisv6 "github.com/ladmakhi81/learnup/pkg/cache/redis/v6"
	"github.com/ladmakhi81/learnup/pkg/env"
	"github.com/ladmakhi81/learnup/pkg/env/koanf"
	ffmpegv1 "github.com/ladmakhi81/learnup/pkg/ffmpeg/v1"
	logrusv1 "github.com/ladmakhi81/learnup/pkg/logger/logrus/v1"
	miniov7 "github.com/ladmakhi81/learnup/pkg/storage/minio/v7"
	jwtv5 "github.com/ladmakhi81/learnup/pkg/token/jwt/v5"
	i18nv2 "github.com/ladmakhi81/learnup/pkg/translations/i18n/v2"
	"github.com/ladmakhi81/learnup/pkg/validation/validator/v10"
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
	courseSvc := courseService.NewCourseServiceImpl(courseRepo, i18nTranslatorSvc, userSvc, categorySvc)
	ffmpegSvc := ffmpegv1.NewFfmpegSvc()
	videoSvc := videoService.NewVideoServiceImpl(minioSvc, ffmpegSvc, logrusSvc, courseSvc, videoRepo, notificationSvc, i18nTranslatorSvc)
	tusHookSvc := tusHookService.NewTusServiceImpl(videoSvc, logrusSvc)

	// middlewares
	middlewares := middleware.NewMiddleware(tokenSvc, redisSvc)

	// handlers
	userAdminHandler := userHandler.NewHandler(userSvc, validationSvc, i18nTranslatorSvc)
	userAuthHandler := authHandler.NewHandler(authSvc, validationSvc, i18nTranslatorSvc)
	categoryAdminHandler := categoryHandler.NewHandler(categorySvc, i18nTranslatorSvc, validationSvc)
	courseAdminHandler := courseHandler.NewHandler(courseSvc, validationSvc, i18nTranslatorSvc, videoSvc)
	videoAdminHandler := videoHandler.NewHandler(validationSvc, videoSvc, i18nTranslatorSvc)
	tusHandler := tusHookHandler.NewTusHookHandler(tusHookSvc)
	notificationAdminHandler := notificationHandler.NewHandler(notificationSvc, i18nTranslatorSvc)

	// modules
	userModule := user.NewModule(userAdminHandler, middlewares)
	authModule := auth.NewModule(userAuthHandler)
	categoryModule := category.NewModule(categoryAdminHandler, middlewares)
	courseModule := course.NewModule(courseAdminHandler, middlewares)
	tusModule := tus.NewModule(tusHandler)
	videoModule := video.NewModule(videoAdminHandler)
	notificationModule := notification.NewModule(notificationAdminHandler, middlewares)

	// register module
	userModule.Register(api)
	authModule.Register(api)
	categoryModule.Register(api)
	courseModule.Register(api)
	tusModule.Register(api)
	videoModule.Register(api)
	notificationModule.Register(api)

	log.Printf("the server running on %s \n", port)

	// run http handler
	log.Fatalln(server.Run(port))
}

func SetupMinio(config *env.EnvConfig) (*minio.Client, error) {
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

func SetupRedis(config *env.EnvConfig) *redis.Client {
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
