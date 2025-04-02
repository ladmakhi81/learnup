package main

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func main() {
	_, minioClientErr := SetupMinio()

	if minioClientErr != nil {
		log.Fatalf("Failed to connect minio: %v", minioClientErr)
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
