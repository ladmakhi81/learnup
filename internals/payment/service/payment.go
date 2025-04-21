package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	paymentDtoReq "github.com/ladmakhi81/learnup/internals/payment/dto/req"
	paymentRepository "github.com/ladmakhi81/learnup/internals/payment/repo"
	transactionDtoReq "github.com/ladmakhi81/learnup/internals/transaction/dto/req"
	transactionService "github.com/ladmakhi81/learnup/internals/transaction/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/types"
)

type PaymentService interface {
	Create(dto paymentDtoReq.CreatePaymentReq) (*entities.Payment, error)
	Verify(dto paymentDtoReq.VerifyPaymentReq) error
	FetchByAuthority(authority string) (*entities.Payment, error)
	FetchPageable(page, pageSize int) ([]*entities.Payment, error)
	FetchCount() (int, error)
}

type PaymentServiceImpl struct {
	zarinpalGateway contracts.PaymentGateway
	zibalGateway    contracts.PaymentGateway
	stripeGateway   contracts.PaymentGateway
	paymentRepo     paymentRepository.PaymentRepo
	config          *dtos.EnvConfig
	translationSvc  contracts.Translator
	transactionSvc  transactionService.TransactionService
}

func NewPaymentService(
	zarinpalGateway contracts.PaymentGateway,
	zibalGateway contracts.PaymentGateway,
	stripeGateway contracts.PaymentGateway,
	paymentRepo paymentRepository.PaymentRepo,
	config *dtos.EnvConfig,
	translationSvc contracts.Translator,
	transactionSvc transactionService.TransactionService,
) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		zarinpalGateway: zarinpalGateway,
		zibalGateway:    zibalGateway,
		stripeGateway:   stripeGateway,
		paymentRepo:     paymentRepo,
		config:          config,
		translationSvc:  translationSvc,
		transactionSvc:  transactionSvc,
	}
}

func (svc PaymentServiceImpl) FetchByAuthority(authority string) (*entities.Payment, error) {
	payment, paymentErr := svc.paymentRepo.FetchByAuthority(authority)
	if paymentErr != nil {
		return nil, types.NewServerError(
			"Error in fetching payment by authority",
			"PaymentServiceImpl.FetchByAuthority",
			paymentErr,
		)
	}
	return payment, nil
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
	if err := svc.paymentRepo.Create(payment); err != nil {
		return nil, err
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
	payment, paymentErr := svc.FetchByAuthority(dto.Authority)
	if paymentErr != nil {
		return paymentErr
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
		transaction, transactionErr := svc.transactionSvc.Create(transactionDtoReq.CreateTransactionReq{
			Currency: svc.getCurrency(dto.Gateway),
			User:     payment.User.FullName(),
			Phone:    payment.User.Phone,
			Amount:   payment.Amount,
			Type:     entities.TransactionType_Deposit,
			Tag:      entities.TransactionTag_Sell,
		})
		if transactionErr != nil {
			return transactionErr
		}
		payment.TransactionID = &transaction.ID
		payment.Status = entities.PaymentStatus_Success
		payment.RefID = resp.RefCode
	} else {
		payment.Status = entities.PaymentStatus_Failure
	}
	if err := svc.paymentRepo.Update(payment); err != nil {
		return types.NewServerError(
			"Error in updating the payment",
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

func (svc PaymentServiceImpl) FetchCount() (int, error) {
	count, countErr := svc.paymentRepo.FetchCount()
	if countErr != nil {
		return 0, countErr
	}
	return count, nil
}

func (svc PaymentServiceImpl) FetchPageable(page, pageSize int) ([]*entities.Payment, error) {
	payments, paymentsErr := svc.paymentRepo.FetchPageable(page, pageSize)
	if paymentsErr != nil {
		return nil, paymentsErr
	}
	return payments, nil
}
