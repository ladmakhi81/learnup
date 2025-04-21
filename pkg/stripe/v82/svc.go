package stripev82

import (
	"fmt"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"net/http"
	"net/url"
)

type StripeClient struct {
	config *dtos.EnvConfig
}

func NewStripeClient(config *dtos.EnvConfig) (*StripeClient, error) {
	stripe.Key = config.Stripe.Key
	proxyURL, err := url.Parse("http://127.0.0.1:12334/")
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
		HTTPClient: client,
	})
	stripe.SetBackend(stripe.APIBackend, backend)
	return &StripeClient{
		config: config,
	}, nil
}

func (svc StripeClient) CreateRequest(dto dtos.CreatePaymentGatewayDto) (*dtos.CreatePaymentGatewayResDto, error) {
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Dynamic Product"),
					},
					UnitAmount: stripe.Int64(int64(dto.Amount)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf(svc.config.Stripe.CallbackURL, "&status=success")),
		CancelURL:  stripe.String(fmt.Sprintf(svc.config.Stripe.CallbackURL, "&status=cancel")),
	}
	req, reqErr := session.New(params)
	if reqErr != nil {
		return nil, reqErr
	}
	return &dtos.CreatePaymentGatewayResDto{
		ID:      req.ID,
		PayLink: req.URL,
	}, nil
}

func (svc StripeClient) VerifyTransaction(dto dtos.VerifyTransactionDto) (*dtos.VerifyTransactionResDto, error) {
	req, reqErr := session.Get(dto.ID, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	isSuccess := req.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid
	return &dtos.VerifyTransactionResDto{
		IsSuccess: isSuccess,
		RefCode:   req.ID,
	}, nil
}
