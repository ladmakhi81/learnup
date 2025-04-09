package storage

import (
	"context"
	"io"
)

type Storage interface {
	BucketExist(context.Context, string) (bool, *StorageError)
	CreateBucket(context.Context, string) *StorageError
	DeleteBucket(context.Context, string) *StorageError
	UploadFileByContent(context.Context, string, string, string, []byte) (*UploadResult, *StorageError)
	GetFile(context.Context, string, string) ([]byte, *StorageError)
	GetFileReader(context.Context, string, string) (io.Reader, *StorageError)
	DeleteObject(context.Context, string, string) error
}
