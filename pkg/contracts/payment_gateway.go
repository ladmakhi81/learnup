package contracts

import "github.com/ladmakhi81/learnup/pkg/dtos"

type PaymentGateway interface {
	CreateRequest(dto dtos.CreatePaymentGatewayDto) (*dtos.CreatePaymentGatewayResDto, error)
	VerifyTransaction(dto dtos.VerifyTransactionDto) (*dtos.VerifyTransactionResDto, error)
}
