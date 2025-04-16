package contracts

import (
	"github.com/ladmakhi81/learnup/pkg/dtos"
)

type HttpClient interface {
	Post(dto dtos.PostRequestDTO) (*dtos.HttpResponse, error)
	Get(dto dtos.GetRequestDTO) (*dtos.HttpResponse, error)
}
