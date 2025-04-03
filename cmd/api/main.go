package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/db"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func main() {
	_, minioClientErr := SetupMinio()
	if minioClientErr != nil {
		log.Fatalf("Failed to connect minio: %v", minioClientErr)
	}

	SetupRedis()

	dbClient := db.NewDatabase()
	if err := dbClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("main function invoked")
}

func SetupMinio() (*minio.Client, error) {
	//TODO: replace these hard code value with config data
	endpoint := "127.0.0.1:9000"
	username := "root_user"
	password := "root_password"
	region := "us-east-1"
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
		Secure: false,
		Region: region,
	})
}

func SetupRedis() *redis.Client {
	//TODO: replace these hard code value with config data
	host := "localhost"
	port := "6379"
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})
}
