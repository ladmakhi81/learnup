package contracts

import (
	"github.com/ladmakhi81/learnup/pkg/dtos"
)

type Mail interface {
	SendPlain(dto dtos.SendMailReq) error
	SendTemplate(dto dtos.SendTemplateMailReq) error
}
