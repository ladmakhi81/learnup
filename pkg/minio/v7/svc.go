package miniov7

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type MinioClientSvc struct {
	minio  *minio.Client
	region string
}

func setupMinioClient(config *dtos.EnvConfig) (*minio.Client, error) {
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

func NewMinioClientSvc(config *dtos.EnvConfig) (*MinioClientSvc, error) {
	client, err := setupMinioClient(config)
	if err != nil {
		return nil, err
	}
	return &MinioClientSvc{minio: client, region: config.Minio.Region}, nil
}

func (svc MinioClientSvc) BucketExist(
	ctx context.Context,
	bucketName string,
) (bool, error) {
	exist, err := svc.minio.BucketExists(ctx, bucketName)
	if err != nil {
		return false, dtos.NewStorageError(
			fmt.Sprintf("Error: bucket exist checking error - %s", err.Error()),
			"MinioClientSvc.BucketExist",
		)
	}
	return exist, nil
}

func (svc MinioClientSvc) CreateBucket(
	ctx context.Context,
	bucketName string,
) error {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return existErr
	}
	if isBucketExist {
		return dtos.NewStorageError(
			"Error: bucket already exists",
			"MinioClientSvc.CreateBucket",
		)
	}
	err := svc.minio.MakeBucket(
		ctx,
		bucketName,
		minio.MakeBucketOptions{Region: svc.region},
	)
	if err != nil {
		return dtos.NewStorageError(
			"Error: happen in create bucket",
			"MinioClientSvc.CreateBucket",
		)
	}
	return nil
}

func (svc MinioClientSvc) DeleteBucket(
	ctx context.Context,
	bucketName string,
) error {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return existErr
	}
	if !isBucketExist {
		return dtos.NewStorageError(
			"Error: bucket does not exist",
			"MinioClientSvc.DeleteBucket",
		)
	}
	err := svc.minio.RemoveBucket(ctx, bucketName)
	if err != nil {
		return dtos.NewStorageError(
			"Error: happen in delete bucket",
			"MinioClientSvc.DeleteBucket",
		)
	}
	return nil
}

func (svc MinioClientSvc) UploadFileByContent(
	ctx context.Context,
	bucketName string,
	objectPath string,
	contentType string,
	fileContents []byte,
) (*dtos.UploadResult, error) {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return nil, existErr
	}
	if !isBucketExist {
		return nil, dtos.NewStorageError(
			"Error: bucket does not exist",
			"MinioClientSvc.UploadFileByContent",
		)
	}
	fileBuffer := bytes.NewBuffer(fileContents)
	info, err := svc.minio.PutObject(
		ctx,
		bucketName,
		objectPath,
		fileBuffer,
		int64(fileBuffer.Len()),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, dtos.NewStorageError(
			"Error: happen in uploading minio file",
			"MinioClientSvc.UploadFileByContent",
		)
	}
	return dtos.NewUploadResult(info.Key, info.Size), nil
}

func (svc MinioClientSvc) GetFile(
	ctx context.Context,
	bucketName,
	fileName string,
) ([]byte, error) {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return nil, existErr
	}
	if !isBucketExist {
		return nil, dtos.NewStorageError(
			"Error: bucket does not exist",
			"MinioClientSvc.GetFile",
		)
	}
	object, objectErr := svc.minio.GetObject(
		ctx,
		bucketName,
		fileName,
		minio.GetObjectOptions{},
	)
	if objectErr != nil {
		return nil, dtos.NewStorageError(
			"Error: happen in get file",
			"MinioClientSvc.GetFile",
		)
	}
	fileContents := make([]byte, 2025)
	contentLength, err := object.Read(fileContents)
	if err != nil {
		return nil, dtos.NewStorageError(
			"Error: happen in reading file contents",
			"MinioClientSvc.GetFile",
		)
	}
	return fileContents[:contentLength], nil
}

func (svc MinioClientSvc) GetFileReader(
	ctx context.Context,
	bucketName,
	fileName string,
) (io.Reader, error) {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return nil, existErr
	}
	if !isBucketExist {
		return nil, dtos.NewStorageError(
			"Error: bucket does not exist",
			"MinioClientSvc.GetFileReader",
		)
	}
	object, objectErr := svc.minio.GetObject(
		ctx,
		bucketName,
		fileName,
		minio.GetObjectOptions{},
	)
	if objectErr != nil {
		return nil, dtos.NewStorageError(
			"Error: happen in get file",
			"MinioClientSvc.GetFileReader",
		)
	}
	return object, nil
}

func (svc MinioClientSvc) DeleteObject(
	ctx context.Context,
	bucketName string,
	objectID string,
) error {
	err := svc.minio.RemoveObject(
		ctx,
		bucketName,
		objectID,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return dtos.NewStorageError(
			"Error: happen in delete object id",
			"MinioClientSvc.DeleteObject",
		)
	}
	return nil
}
