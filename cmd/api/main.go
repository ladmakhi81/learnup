package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/auth"
	authHandler "github.com/ladmakhi81/learnup/internals/auth/handler"
	authService "github.com/ladmakhi81/learnup/internals/auth/service"
	"github.com/ladmakhi81/learnup/internals/user"
	userHandler "github.com/ladmakhi81/learnup/internals/user/handler"
	"github.com/ladmakhi81/learnup/internals/user/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	redisv6 "github.com/ladmakhi81/learnup/pkg/cache/redis/v6"
	"github.com/ladmakhi81/learnup/pkg/env"
	"github.com/ladmakhi81/learnup/pkg/env/koanf"
	jwtv5 "github.com/ladmakhi81/learnup/pkg/token/jwt/v5"
	"github.com/ladmakhi81/learnup/pkg/validation/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"

	_ "github.com/ladmakhi81/learnup/docs"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Learnup
// @version         1.0
// @BasePath  /api
func main() {
	// config file loader
	koanfConfigProvider := koanf.NewKoanfEnvSvc()
	config, configErr := koanfConfigProvider.LoadLearnUp()
	if configErr != nil {
		log.Fatalf("load learn up config failed: %v", configErr)
	}

	// minio
	_, minioClientErr := SetupMinio(config)
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

	// swagger documentation
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// repos
	userRepo := repo.NewUserRepoImpl(dbClient)

	// svcs
	redisSvc := redisv6.NewRedisClientSvc(redisClient)
	tokenSvc := jwtv5.NewJwtSvc(config)
	validationSvc := validatorv10.NewValidatorSvc(validator.New())
	userSvc := userService.NewUserSvcImpl(userRepo)
	authSvc := authService.NewAuthServiceImpl(userSvc, redisSvc, tokenSvc)

	// handlers
	userAdminHandler := userHandler.NewUserAdminHandler(userSvc, validationSvc)
	userAuthHandler := authHandler.NewUserAuthHandler(authSvc, validationSvc)

	// modules
	userModule := user.NewModule(userAdminHandler)
	authModule := auth.NewAuthModule(userAuthHandler)

	// register module
	userModule.Register(api)
	authModule.Register(api)

	log.Printf("the server running on %s \n", port)

	// run http handler
	log.Fatalln(server.Run(port))
}

func SetupMinio(config *env.EnvConfig) (*minio.Client, error) {
	endpoint := config.Minio.URL
	username := config.Minio.Username
	password := config.Minio.Password
	region := config.Minio.Region
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
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
