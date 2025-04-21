package dtoreq

import "github.com/ladmakhi81/learnup/db/entities"

type VerifyPaymentReq struct {
	Authority string
	Gateway   entities.PaymentGateway
}
