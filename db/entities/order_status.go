package entities

type OrderStatus string

const (
	OrderStatus_Pending OrderStatus = "pending"
	OrderStatus_Success OrderStatus = "success"
	OrderStatus_Failed  OrderStatus = "failed"
)
