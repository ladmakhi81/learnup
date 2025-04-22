package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	paymentDtoReq "github.com/ladmakhi81/learnup/internals/payment/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/types"
	"time"
)

type PaymentService interface {
	Create(dto paymentDtoReq.CreatePaymentReq) (*entities.Payment, error)
	Verify(dto paymentDtoReq.VerifyPaymentReq) error
	FetchPageable(page, pageSize int) ([]*entities.Payment, int, error)
}

type PaymentServiceImpl struct {
	repo            *db.Repositories
	zarinpalGateway contracts.PaymentGateway
	zibalGateway    contracts.PaymentGateway
	stripeGateway   contracts.PaymentGateway
	config          *dtos.EnvConfig
	translationSvc  contracts.Translator
}

func NewPaymentService(
	repo *db.Repositories,
	zarinpalGateway contracts.PaymentGateway,
	zibalGateway contracts.PaymentGateway,
	stripeGateway contracts.PaymentGateway,
	config *dtos.EnvConfig,
	translationSvc contracts.Translator,
) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		zarinpalGateway: zarinpalGateway,
		zibalGateway:    zibalGateway,
		stripeGateway:   stripeGateway,
		config:          config,
		translationSvc:  translationSvc,
		repo:            repo,
	}
}

func (svc PaymentServiceImpl) Create(dto paymentDtoReq.CreatePaymentReq) (*entities.Payment, error) {
	gateway := svc.selectGateway(dto.Gateway)
	if gateway == nil {
		return nil, types.NewBadRequestError(
			svc.translationSvc.Translate("payment.errors.gateway_not_found"),
		)
	}
	merchant := svc.getMerchantID(dto.Gateway)
	if merchant == "" {
		return nil, types.NewBadRequestError(
			svc.translationSvc.Translate("payment.errors.merchant_not_found"),
		)
	}
	resp, respErr := gateway.CreateRequest(dtos.CreatePaymentGatewayDto{
		Amount: dto.Amount,
	})
	if respErr != nil {
		return nil, types.NewServerError(
			"Error in generating URL for gateway and initialize request",
			"PaymentServiceImpl.Create",
			respErr,
		)
	}
	payment := &entities.Payment{
		Amount:     dto.Amount,
		Authority:  resp.ID,
		MerchantID: merchant,
		Gateway:    dto.Gateway,
		Status:     entities.PaymentStatus_Pending,
		UserID:     dto.UserID,
		OrderID:    dto.OrderID,
		PayLink:    resp.PayLink,
	}
	if err := svc.repo.PaymentRepo.Create(payment); err != nil {
		return nil, types.NewServerError(
			"Error in creating payment",
			"PaymentServiceImpl.Create",
			err,
		)
	}
	return payment, nil
}

func (svc PaymentServiceImpl) Verify(dto paymentDtoReq.VerifyPaymentReq) error {
	gateway := svc.selectGateway(dto.Gateway)
	if gateway == nil {
		return types.NewBadRequestError(
			svc.translationSvc.Translate("payment.errors.gateway_not_found"),
		)
	}
	payment, paymentErr := svc.repo.PaymentRepo.GetOne(map[string]any{
		"authority": dto.Authority,
	}, []string{"User", "Order"})
	if paymentErr != nil {
		return types.NewServerError(
			"Error in fetching payment by authority",
			"PaymentServiceImpl.Verify",
			paymentErr,
		)
	}
	if payment == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("payment.errors.not_found"),
		)
	}
	resp, respErr := gateway.VerifyTransaction(dtos.VerifyTransactionDto{
		ID:     payment.Authority,
		Amount: payment.Amount,
	})
	if respErr != nil {
		return types.NewServerError(
			"Error in verifying transaction from server",
			"PaymentServiceImpl.Verify",
			respErr,
		)
	}
	if resp.IsSuccess {
		transaction := &entities.Transaction{
			Amount:   payment.Amount,
			Type:     entities.TransactionType_Deposit,
			User:     payment.User.FullName(),
			Phone:    payment.User.Phone,
			Tag:      entities.TransactionTag_Sell,
			Currency: svc.getCurrency(dto.Gateway),
		}
		if err := svc.repo.TransactionRepo.Create(transaction); err != nil {
			return types.NewServerError(
				"Error in creating transaction based on the payment",
				"PaymentServiceImpl.Verify",
				err,
			)
		}
		payment.TransactionID = &transaction.ID
		payment.Status = entities.PaymentStatus_Success
		payment.RefID = resp.RefCode
		payment.Order.Status = entities.OrderStatus_Success

	} else {
		payment.Status = entities.PaymentStatus_Failure
		payment.Order.Status = entities.OrderStatus_Failed
	}
	if err := svc.repo.PaymentRepo.Update(payment); err != nil {
		return types.NewServerError(
			"Error in updating the payment",
			"PaymentServiceImpl.Verify",
			err,
		)
	}
	now := time.Now()
	payment.Order.StatusChangedAt = &now
	if err := svc.repo.OrderRepo.Update(payment.Order); err != nil {
		return types.NewServerError(
			"Error in updating status of order",
			"PaymentServiceImpl.Verify",
			err,
		)
	}
	return nil
}

func (svc PaymentServiceImpl) selectGateway(paymentGateway entities.PaymentGateway) contracts.PaymentGateway {
	switch paymentGateway {
	case entities.PaymentGateway_Zarinpal:
		return svc.zarinpalGateway
	case entities.PaymentGateway_Zibal:
		return svc.zibalGateway
	case entities.PaymentGateway_Stripe:
		return svc.stripeGateway
	default:
		return nil
	}
}

func (svc PaymentServiceImpl) getMerchantID(paymentGateway entities.PaymentGateway) string {
	switch paymentGateway {
	case entities.PaymentGateway_Zarinpal:
		return svc.config.Zarinpal.Merchant
	case entities.PaymentGateway_Zibal:
		return svc.config.Zibal.Merchant
	case entities.PaymentGateway_Stripe:
		return svc.config.Stripe.Key
	default:
		return ""
	}
}

func (svc PaymentServiceImpl) getCurrency(paymentGateway entities.PaymentGateway) string {
	switch paymentGateway {
	case entities.PaymentGateway_Zarinpal:
		return "IRT"
	case entities.PaymentGateway_Zibal:
		return "IRT"
	case entities.PaymentGateway_Stripe:
		return "USD"
	default:
		return ""
	}
}

func (svc PaymentServiceImpl) FetchPageable(page, pageSize int) ([]*entities.Payment, int, error) {
	payments, count, paymentsErr := svc.repo.PaymentRepo.GetPaginated(
		repositories.GetPaginatedOptions{
			Offset: &page,
			Limit:  &pageSize,
		},
	)
	if paymentsErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching payments pages",
			"PaymentServiceImpl.FetchPageable",
			paymentsErr,
		)
	}
	return payments, count, nil
}
