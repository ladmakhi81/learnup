package dtos

type StorageError struct {
	Message  string
	Location string
}

func (e StorageError) Error() string {
	return e.Message
}

func NewStorageError(message string, location string) *StorageError {
	return &StorageError{
		Message:  message,
		Location: location,
	}
}

type UploadResult struct {
	ObjectID string
	Size     int64
}

func NewUploadResult(objectID string, size int64) *UploadResult {
	return &UploadResult{
		ObjectID: objectID,
		Size:     size,
	}
}
