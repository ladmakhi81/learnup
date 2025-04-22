package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	paymentDtoReq "github.com/ladmakhi81/learnup/internals/payment/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/types"
	"log"
	"time"
)

type PaymentService interface {
	Create(tx db.UnitOfWorkTx, dto paymentDtoReq.CreatePaymentReq) (*entities.Payment, error)
	Verify(dto paymentDtoReq.VerifyPaymentReq) error
	FetchPageable(page, pageSize int) ([]*entities.Payment, int, error)
}

type PaymentServiceImpl struct {
	unitOfWork      db.UnitOfWork
	zarinpalGateway contracts.PaymentGateway
	zibalGateway    contracts.PaymentGateway
	stripeGateway   contracts.PaymentGateway
	config          *dtos.EnvConfig
	translationSvc  contracts.Translator
}

func NewPaymentService(
	unitOfWork db.UnitOfWork,
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
		unitOfWork:      unitOfWork,
	}
}

func (svc PaymentServiceImpl) Create(tx db.UnitOfWorkTx, dto paymentDtoReq.CreatePaymentReq) (*entities.Payment, error) {
	const operationName = "PaymentServiceImpl.Create"
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
	resp, err := gateway.CreateRequest(dtos.CreatePaymentGatewayDto{
		Amount: dto.Amount,
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in generating URL for gateway and initialize request",
			operationName,
			err,
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
	if err := tx.PaymentRepo().Create(payment); err != nil {
		return nil, types.NewServerError(
			"Error in creating payment",
			operationName,
			err,
		)
	}
	return payment, nil
}

func (svc PaymentServiceImpl) Verify(dto paymentDtoReq.VerifyPaymentReq) error {
	const operationName = "PaymentServiceImpl.Verify"
	transactionID, err := db.WithTx(svc.unitOfWork, func(tx db.UnitOfWorkTx) (uint, error) {
		gateway := svc.selectGateway(dto.Gateway)
		if gateway == nil {
			return 0, types.NewBadRequestError(
				svc.translationSvc.Translate("payment.errors.gateway_not_found"),
			)
		}
		payment, err := tx.PaymentRepo().GetOne(map[string]any{
			"authority": dto.Authority,
		}, []string{"User", "Order"})
		if err != nil {
			return 0, types.NewServerError(
				"Error in fetching payment by authority",
				operationName,
				err,
			)
		}
		if payment == nil {
			return 0, types.NewNotFoundError(
				svc.translationSvc.Translate("payment.errors.not_found"),
			)
		}
		resp, err := gateway.VerifyTransaction(dtos.VerifyTransactionDto{
			ID:     payment.Authority,
			Amount: payment.Amount,
		})
		if err != nil {
			return 0, types.NewServerError(
				"Error in verifying transaction from server",
				operationName,
				err,
			)
		}
		var transactionID uint
		if resp.IsSuccess {
			transaction := &entities.Transaction{
				Amount:   payment.Amount,
				Type:     entities.TransactionType_Deposit,
				User:     payment.User.FullName(),
				Phone:    payment.User.Phone,
				Tag:      entities.TransactionTag_Sell,
				Currency: svc.getCurrency(dto.Gateway),
			}
			if err := tx.TransactionRepo().Create(transaction); err != nil {
				return 0, types.NewServerError(
					"Error in creating transaction based on the payment",
					operationName,
					err,
				)
			}
			transactionID = transaction.ID
			payment.TransactionID = &transaction.ID
			payment.Status = entities.PaymentStatus_Success
			payment.RefID = resp.RefCode
			payment.Order.Status = entities.OrderStatus_Success

		} else {
			payment.Status = entities.PaymentStatus_Failure
			payment.Order.Status = entities.OrderStatus_Failed
		}
		if err := tx.PaymentRepo().Update(payment); err != nil {
			return 0, types.NewServerError(
				"Error in updating the payment",
				operationName,
				err,
			)
		}
		now := time.Now()
		payment.Order.StatusChangedAt = &now
		if err := tx.OrderRepo().Update(payment.Order); err != nil {
			return 0, types.NewServerError(
				"Error in updating status of order",
				operationName,
				err,
			)
		}
		return transactionID, nil
	})
	if err != nil {
		return err
	}
	log.Printf("New Transaction Generated : %v \n", transactionID)
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
	const operationName = "PaymentServiceImpl.FetchPageable"
	payments, count, err := svc.unitOfWork.PaymentRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset: &page,
			Limit:  &pageSize,
		},
	)
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching payments pages",
			operationName,
			err,
		)
	}
	return payments, count, nil
}
