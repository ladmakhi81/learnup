package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/pkg/env"
	"github.com/ladmakhi81/learnup/pkg/env/koanf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

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
	SetupRedis(config)

	// database
	dbClient := db.NewDatabase(config)
	if err := dbClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("main function invoked")
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
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})
}
