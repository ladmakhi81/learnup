package dtoreq

type VerifyCourseReq struct {
	ID                          uint    `json:"-"`
	Fee                         float64 `json:"fee" validate:"omitempty,gte=0"`
	DiscountFeeAmountPercentage float64 `json:"discountFeeAmountPercentage" validate:"omitempty,gte=0"`
}
