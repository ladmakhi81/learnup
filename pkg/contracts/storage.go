package contracts

import (
	"context"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"io"
)

type Storage interface {
	BucketExist(ctx context.Context, bucketName string) (bool, *dtos.StorageError)
	CreateBucket(ctx context.Context, bucketName string) *dtos.StorageError
	DeleteBucket(ctx context.Context, bucketName string) *dtos.StorageError
	UploadFileByContent(ctx context.Context, bucketName string, objectPath string, contentType string, fileContents []byte) (*dtos.UploadResult, *dtos.StorageError)
	GetFile(ctx context.Context, bucketName string, fileName string) ([]byte, *dtos.StorageError)
	GetFileReader(ctx context.Context, bucketName string, fileName string) (io.Reader, *dtos.StorageError)
	DeleteObject(ctx context.Context, bucketName string, objectId string) error
}
