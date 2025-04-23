package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type GetPageablePaymentItemDto struct {
	ID        uint                     `json:"id"`
	CreatedAt time.Time                `json:"createdAt"`
	UpdatedAt time.Time                `json:"updatedAt"`
	UserID    uint                     `json:"userId"`
	Gateway   entities2.PaymentGateway `json:"gateway"`
	Status    entities2.PaymentStatus  `json:"status"`
	Amount    float64                  `json:"amount"`
}

func MapGetPageablePaymentItemsDto(payments []*entities2.Payment) []*GetPageablePaymentItemDto {
	res := make([]*GetPageablePaymentItemDto, len(payments))
	for i, payment := range payments {
		res[i] = &GetPageablePaymentItemDto{
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
