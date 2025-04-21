package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"time"
)

type GetPageablePaymentItem struct {
	ID        uint                    `json:"id"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
	UserID    uint                    `json:"userId"`
	Gateway   entities.PaymentGateway `json:"gateway"`
	Status    entities.PaymentStatus  `json:"status"`
	Amount    float64                 `json:"amount"`
}

func MapGetPageablePaymentItems(payments []*entities.Payment) []*GetPageablePaymentItem {
	res := make([]*GetPageablePaymentItem, len(payments))
	for i, payment := range payments {
		res[i] = &GetPageablePaymentItem{
			ID:        payment.ID,
			CreatedAt: payment.CreatedAt,
			UpdatedAt: payment.UpdatedAt,
			UserID:    payment.UserID,
			Gateway:   payment.Gateway,
			Status:    payment.Status,
			Amount:    payment.Amount,
		}
	}
	return res
}
