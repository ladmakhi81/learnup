package entities

type PaymentGateway string

const (
	PaymentGateway_Zibal    = "zibal"
	PaymentGateway_Zarinpal = "zarinpal"
	PaymentGateway_Stripe   = "stripe"
)
