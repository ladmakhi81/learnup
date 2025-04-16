package restyv2

import (
	"github.com/go-resty/resty/v2"
	"github.com/ladmakhi81/learnup/pkg/dtos"
)

type RestyHttpSvc struct {
	httpClient *resty.Client
}

func NewRestyHttpSvc() *RestyHttpSvc {
	return &RestyHttpSvc{
		httpClient: resty.New(),
	}
}

func (svc RestyHttpSvc) Post(dto dtos.PostRequestDTO) (*dtos.HttpResponse, error) {
	resp, respErr := svc.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(dto).
		Post(dto.URL)
	if respErr != nil {
		return nil, dtos.NewHttpError(
			"Error: happen in sending post request",
			"RestyHttpSvc.Post",
		)
	}
	return dtos.NewHttpResponse(
		resp.StatusCode(),
		resp.Body(),
	), nil
}

func (svc RestyHttpSvc) Get(dto dtos.GetRequestDTO) (*dtos.HttpResponse, error) {
	resp, respErr := svc.httpClient.R().
		SetQueryParams(dto.QueryParams).
		SetHeader("Content-Type", "application/json").
		Get(dto.URL)
	if respErr != nil {
		return nil, dtos.NewHttpError(
			"Error: happen in sending get request",
			"RestyHttpSvc.Get",
		)
	}
	return dtos.NewHttpResponse(
		resp.StatusCode(),
		resp.Body(),
	), nil
}
