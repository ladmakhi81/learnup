package zibalv1

import "time"

type CreateRequestDTO struct {
	Merchant    string  `json:"merchant"`
	Amount      float64 `json:"amount"`
	CallbackURL string  `json:"callbackUrl"`
}

type CreateRequestResDTO struct {
	Message string `json:"message"`
	Result  uint   `json:"result"`
	TrackID int    `json:"trackId"`
}

type VerifyRequestDTO struct {
	Merchant string `json:"merchant"`
	TrackID  int    `json:"trackId"`
}

type VerifyRequestResDTO struct {
	Message     string    `json:"message"`
	Result      uint      `json:"result"`
	RefNumber   string    `json:"refNumber"`
	PaidAt      time.Time `json:"paidAt"`
	Status      uint      `json:"status"`
	Amount      float64   `json:"amount"`
	OrderID     string    `json:"orderId"`
	Description string    `json:"description"`
	CardNumber  string    `json:"cardNumber"`
}
