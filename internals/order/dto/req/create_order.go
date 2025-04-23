package dtoreq

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
)

type CreateOrderReq struct {
	Carts   []uint                  `json:"carts" validate:"required,dive,gte=1"`
	Gateway entities.PaymentGateway `json:"gateway" validate:"required,oneof=zibal zarinpal stripe"`
}
