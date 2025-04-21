package zarinpalv1

type CreateRequestDTO struct {
	MerchantID  string  `json:"merchant_id"`
	Amount      float64 `json:"amount"`
	CallbackURL string  `json:"callback_url"`
	Description string  `json:"description"`
}

type CreateRequestResDTO struct {
	Data struct {
		Authority string `json:"authority"`
		Fee       uint   `json:"fee"`
		FeeType   string `json:"fee_type"`
		Code      uint   `json:"code"`
		Message   string `json:"message"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
		Code    uint   `json:"code"`
	} `json:"errors"`
}

type VerifyRequestDTO struct {
	MerchantID string  `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	Authority  string  `json:"authority"`
}

type VerifyRequestResDTO struct {
	Data struct {
		Wages       *string `json:"wages"`
		Code        int     `json:"code"`
		Message     string  `json:"message"`
		CardHash    string  `json:"card_hash"`
		CardPan     string  `json:"card_pan"`
		RefID       int     `json:"ref_id"`
		FeeType     string  `json:"fee_type"`
		Fee         int     `json:"fee"`
		ShaparakFee int     `json:"shaparak_fee"`
		OrderID     *string `json:"order_id"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
		Code    uint   `json:"code"`
	} `json:"errors"`
}
