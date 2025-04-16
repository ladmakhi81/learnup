package dtoreq

type VerifyCourseReq struct {
	ID                          uint    `json:"-"`
	Fee                         float64 `json:"fee" validate:"required,gte=0"`
	DiscountFeeAmountPercentage float64 `json:"discountFeeAmountPercentage" validate:"required,gte=0"`
}
