package storage

import (
	"context"
	"io"
)

type Storage interface {
	BucketExist(ctx context.Context, bucketName string) (bool, *StorageError)
	CreateBucket(ctx context.Context, bucketName string) *StorageError
	DeleteBucket(ctx context.Context, bucketName string) *StorageError
	UploadFileByContent(ctx context.Context, bucketName string, objectPath string, contentType string, fileContents []byte) (*UploadResult, *StorageError)
	GetFile(ctx context.Context, bucketName string, fileName string) ([]byte, *StorageError)
	GetFileReader(ctx context.Context, bucketName string, fileName string) (io.Reader, *StorageError)
	DeleteObject(ctx context.Context, bucketName string, objectId string) error
}
