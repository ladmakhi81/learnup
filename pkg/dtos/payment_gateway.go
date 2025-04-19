package dtos

type CreatePaymentGatewayDto struct {
	CallbackURL string
	Amount      float64
}

type CreatePaymentGatewayResDto struct {
	PayLink string
	ID      string
}

type VerifyTransactionDto struct {
	ID     string
	Amount float64
}

type VerifyTransactionResDto struct {
	IsSuccess bool
	RefCode   string
}
