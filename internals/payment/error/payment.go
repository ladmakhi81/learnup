package error

import "github.com/ladmakhi81/learnup/types"

var (
	Payment_GatewayNotFound  = types.NewBadRequestError("payment.errors.gateway_not_found")
	Payment_MerchantNotFound = types.NewBadRequestError("payment.errors.merchant_not_found")
	Payment_NotFound         = types.NewNotFoundError("payment.errors.not_found")
)
