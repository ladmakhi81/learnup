package entities

type TransactionType string

const (
	TransactionType_Withdraw TransactionType = "withdraw"
	TransactionType_Deposit  TransactionType = "deposit"
)
