package http

type HttpError struct {
	Message  string
	Location string
}

func (e HttpError) Error() string {
	return e.Message
}

func NewHttpError(
	message string,
	location string,
) *HttpError {
	return &HttpError{
		Message:  message,
		Location: location,
	}
}

type HttpResponse struct {
	StatusCode int
	Result     []byte
}

func NewHttpResponse(statusCode int, result []byte) *HttpResponse {
	return &HttpResponse{
		StatusCode: statusCode,
		Result:     result,
	}
}

// ------------------------------------------------------
// Http Post Function
type PostRequestDTO struct {
	URL  string
	Body any
}

func NewPostRequestDTO(url string, body any) PostRequestDTO {
	return PostRequestDTO{
		URL:  url,
		Body: body,
	}
}

// ------------------------------------------------------
// Http Get Function
type GetRequestDTO struct {
	URL         string
	QueryParams map[string]string
}

func NewGetRequestDTO(url string, queryParams map[string]string) GetRequestDTO {
	return GetRequestDTO{
		URL:         url,
		QueryParams: queryParams,
	}
}
