package types

type ApiResponse struct {
	StatusCode int `json:"statusCode"`
	Data       any `json:"data"`
}

func NewApiResponse(statusCode int, data any) *ApiResponse {
	return &ApiResponse{
		StatusCode: statusCode,
		Data:       data,
	}
}

type ApiError struct {
	StatusCode int    `json:"statusCode"`
	Message    any    `json:"message"`
	Timestamp  int64  `json:"timestamp"`
	TraceID    string `json:"traceId"`
}

func NewApiError(statusCode int, message any, timestamp int64, traceID string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Message:    message,
		Timestamp:  timestamp,
		TraceID:    traceID,
	}
}
