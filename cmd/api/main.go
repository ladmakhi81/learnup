package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/pkg/env/koanf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func main() {
	// config file loader
	koanfConfigProvider := koanf.NewKoanfEnvSvc()
	_, configErr := koanfConfigProvider.LoadLearnUp()
	if configErr != nil {
		log.Fatalf("load learn up config failed: %v", configErr)
	}

	// minio
	_, minioClientErr := SetupMinio()
	if minioClientErr != nil {
		log.Fatalf("Failed to connect minio: %v", minioClientErr)
	}

	// redis
	SetupRedis()

	// database
	dbClient := db.NewDatabase()
	if err := dbClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("main function invoked")
}

func SetupMinio() (*minio.Client, error) {
	//TODO: replace these hard code value with config data
	endpoint := "learnup_minio_storage_service:9000"
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
	host := "learnup_redis_service"
	port := "6379"
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})
}

//
//func LoadEnvConfig() error {
//	k := koanf.New(".")
//	provider := env.Provider("LEARNUP_", "__", func(s string) string {
//		return strings.ToLower(strings.TrimPrefix(s, "LEARNUP_"))
//	})
//	if err := k.Load(provider, nil); err != nil {
//		return err
//	}
//	var envData struct {
//		MAIN_DB struct {
//			HOST string `koanf:"host"`
//		} `koanf:"main_db"`
//	}
//	if err := k.Unmarshal("", &envData); err != nil {
//		return err
//	}
//	fmt.Println("data: ", envData.MAIN_DB.HOST)
//	return nil
//}
