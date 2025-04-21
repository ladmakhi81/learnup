package dtoreq

import "github.com/ladmakhi81/learnup/db/entities"

type CreatePaymentReq struct {
	OrderID uint
	UserID  uint
	Gateway entities.PaymentGateway
	Amount  float64
}
