package storage

import "context"

type Storage interface {
	BucketExist(context.Context, string) (bool, *StorageError)
	CreateBucket(context.Context, string) *StorageError
	DeleteBucket(context.Context, string) *StorageError
	UploadFileByContent(context.Context, string, []byte) (*UploadResult, *StorageError)
	GetFile(context.Context, string, string) ([]byte, *StorageError)
}
