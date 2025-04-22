package contracts

import (
	"context"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"io"
)

type Storage interface {
	BucketExist(ctx context.Context, bucketName string) (bool, error)
	CreateBucket(ctx context.Context, bucketName string) error
	DeleteBucket(ctx context.Context, bucketName string) error
	UploadFileByContent(ctx context.Context, bucketName string, objectPath string, contentType string, fileContents []byte) (*dtos.UploadResult, error)
	GetFile(ctx context.Context, bucketName string, fileName string) ([]byte, error)
	GetFileReader(ctx context.Context, bucketName string, fileName string) (io.Reader, error)
	DeleteObject(ctx context.Context, bucketName string, objectId string) error
}
