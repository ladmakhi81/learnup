package dtoreq

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
)

type VerifyPaymentReqDto struct {
	Authority string
	Gateway   entities.PaymentGateway
}
