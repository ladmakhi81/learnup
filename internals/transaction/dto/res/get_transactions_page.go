package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type GetTransactionPageableItem struct {
	ID        uint                      `json:"id"`
	CreatedAt time.Time                 `json:"createdAt"`
	Amount    float64                   `json:"amount"`
	User      string                    `json:"user"`
	Phone     string                    `json:"phone"`
	Type      entities2.TransactionType `json:"type"`
	Tag       entities2.TransactionTag  `json:"tag"`
	Currency  string                    `json:"currency"`
}

func MapGetTransactionPageableItems(transactions []*entities2.Transaction) []*GetTransactionPageableItem {
	res := make([]*GetTransactionPageableItem, len(transactions))
	for i, item := range transactions {
		res[i] = &GetTransactionPageableItem{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			Amount:    item.Amount,
			User:      item.User,
			Phone:     item.Phone,
			Type:      item.Type,
			Tag:       item.Tag,
			Currency:  item.Currency,
		}
	}
	return res
}
