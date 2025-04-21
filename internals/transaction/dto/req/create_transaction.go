package dtoreq

import "github.com/ladmakhi81/learnup/db/entities"

type CreateTransactionReq struct {
	Amount   float64
	User     string
	Phone    string
	Type     entities.TransactionType
	Tag      entities.TransactionTag
	Currency string
}
