package http

type HttpClient interface {
	Post(dto PostRequestDTO) (*HttpResponse, error)
	Get(dto GetRequestDTO) (*HttpResponse, error)
}
