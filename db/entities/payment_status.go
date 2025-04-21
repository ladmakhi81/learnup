package entities

type PaymentStatus string

const (
	PaymentStatus_Pending PaymentStatus = "pending"
	PaymentStatus_Success PaymentStatus = "success"
	PaymentStatus_Failure PaymentStatus = "failed"
)
