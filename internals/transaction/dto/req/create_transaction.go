package dtoreq

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
)

type CreateTransactionReq struct {
	Amount   float64
	User     string
	Phone    string
	Type     entities2.TransactionType
	Tag      entities2.TransactionTag
	Currency string
}
