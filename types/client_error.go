package types

import "net/http"

type ClientError struct {
	StatusCode int
	Message    string
	Metadata   any `json:"errors"`
}

func (e ClientError) Error() string {
	return e.Message
}

func NewConflictError(message string) *ClientError {
	return &ClientError{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}

func NewBadRequestError(message string) *ClientError {
	return &ClientError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewBadRequestDTOError(metadata any) *ClientError {
	return &ClientError{
		StatusCode: http.StatusBadRequest,
		Metadata:   metadata,
	}
}
