package miniov7

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/ladmakhi81/learnup/pkg/storage"
	"github.com/minio/minio-go/v7"
)

type MinioClientSvc struct {
	minio  *minio.Client
	region string
}

func NewMinioClientSvc(
	minio *minio.Client,
	region string,
) *MinioClientSvc {
	return &MinioClientSvc{
		minio:  minio,
		region: region,
	}
}

func (svc MinioClientSvc) BucketExist(
	ctx context.Context,
	bucketName string,
) (bool, *storage.StorageError) {
	exist, err := svc.minio.BucketExists(ctx, bucketName)
	if err != nil {
		return false, storage.NewStorageError(
			"Error: bucket exist checking error",
			"MinioClientSvc.BucketExist",
		)
	}
	return exist, nil
}

func (svc MinioClientSvc) CreateBucket(
	ctx context.Context,
	bucketName string,
) *storage.StorageError {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return existErr
	}
	if isBucketExist {
		return storage.NewStorageError(
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
		return storage.NewStorageError(
			"Error: happen in create bucket",
			"MinioClientSvc.CreateBucket",
		)
	}
	return nil
}

func (svc MinioClientSvc) DeleteBucket(
	ctx context.Context,
	bucketName string,
) *storage.StorageError {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return existErr
	}
	if !isBucketExist {
		return storage.NewStorageError(
			"Error: bucket does not exist",
			"MinioClientSvc.DeleteBucket",
		)
	}
	err := svc.minio.RemoveBucket(ctx, bucketName)
	if err != nil {
		return storage.NewStorageError(
			"Error: happen in delete bucket",
			"MinioClientSvc.DeleteBucket",
		)
	}
	return nil
}

func (svc MinioClientSvc) UploadFileByContent(
	ctx context.Context,
	bucketName string,
	contentType string,
	fileContents []byte,
) (*storage.UploadResult, *storage.StorageError) {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return nil, existErr
	}
	if !isBucketExist {
		return nil, storage.NewStorageError(
			"Error: bucket does not exist",
			"MinioClientSvc.UploadFileByContent",
		)
	}
	objectID, objectIDErr := uuid.NewUUID()
	if objectIDErr != nil {
		return nil, storage.NewStorageError(
			"Error: happen in generating object id",
			"MinioClientSvc.UploadFileByContent",
		)
	}
	fileBuffer := bytes.NewBuffer(fileContents)
	info, err := svc.minio.PutObject(
		ctx,
		bucketName,
		objectID.String(),
		fileBuffer,
		int64(fileBuffer.Len()),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, storage.NewStorageError(
			"Error: happen in uploading minio file",
			"MinioClientSvc.UploadFileByContent",
		)
	}
	return storage.NewUploadResult(info.Key, info.Size), nil
}

func (svc MinioClientSvc) GetFile(
	ctx context.Context,
	bucketName,
	fileName string) ([]byte, *storage.StorageError) {
	isBucketExist, existErr := svc.BucketExist(ctx, bucketName)
	if existErr != nil {
		return nil, existErr
	}
	if !isBucketExist {
		return nil, storage.NewStorageError(
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
		return nil, storage.NewStorageError(
			"Error: happen in get file",
			"MinioClientSvc.GetFile",
		)
	}
	fileContents := make([]byte, 2025)
	contentLength, err := object.Read(fileContents)
	if err != nil {
		return nil, storage.NewStorageError(
			"Error: happen in reading file contents",
			"MinioClientSvc.GetFile",
		)
	}
	return fileContents[:contentLength], nil
}
