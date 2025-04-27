package service

import (
	paymentDtoReq "github.com/ladmakhi81/learnup/internals/payment/dto/req"
	paymentError "github.com/ladmakhi81/learnup/internals/payment/error"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/db/repositories"
	"github.com/ladmakhi81/learnup/shared/types"
	"github.com/ladmakhi81/learnup/shared/utils"
	"gorm.io/gorm"
	"log"
)

type PaymentService interface {
	Create(tx db.UnitOfWorkTx, dto paymentDtoReq.CreatePaymentReqDto) (*entities.Payment, error)
	Verify(dto paymentDtoReq.VerifyPaymentReqDto) error
	FetchPageable(page, pageSize int) ([]*entities.Payment, int, error)
}

type paymentService struct {
	unitOfWork      db.UnitOfWork
	zarinpalGateway contracts.PaymentGateway
	zibalGateway    contracts.PaymentGateway
	stripeGateway   contracts.PaymentGateway
	config          *dtos.EnvConfig
}

func NewPaymentService(
	unitOfWork db.UnitOfWork,
	zarinpalGateway contracts.PaymentGateway,
	zibalGateway contracts.PaymentGateway,
	stripeGateway contracts.PaymentGateway,
	config *dtos.EnvConfig,
) PaymentService {
	return &paymentService{
		zarinpalGateway: zarinpalGateway,
		zibalGateway:    zibalGateway,
		stripeGateway:   stripeGateway,
		config:          config,
		unitOfWork:      unitOfWork,
	}
}

func (svc paymentService) Create(tx db.UnitOfWorkTx, dto paymentDtoReq.CreatePaymentReqDto) (*entities.Payment, error) {
	const operationName = "paymentService.Create"
	gateway := svc.selectGateway(dto.Gateway)
	if gateway == nil {
		return nil, paymentError.Payment_GatewayNotFound
	}
	merchant := svc.getMerchantID(dto.Gateway)
	if merchant == "" {
		return nil, paymentError.Payment_MerchantNotFound
	}
	resp, err := gateway.CreateRequest(dtos.CreatePaymentGatewayDto{Amount: dto.Amount})
	if err != nil {
		return nil, types.NewServerError("Error in generating URL for gateway and initialize request", operationName, err)
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
		return nil, types.NewServerError("Error in creating payment", operationName, err)
	}
	return payment, nil
}

func (svc paymentService) Verify(dto paymentDtoReq.VerifyPaymentReqDto) error {
	const operationName = "paymentService.Verify"
	transactionID, err := db.WithTx(svc.unitOfWork, func(tx db.UnitOfWorkTx) (uint, error) {
		gateway := svc.selectGateway(dto.Gateway)
		if gateway == nil {
			return 0, paymentError.Payment_GatewayNotFound
		}
		payment, err := tx.PaymentRepo().GetOne(map[string]any{"authority": dto.Authority}, []string{"User"})
		if err != nil {
			return 0, types.NewServerError("Error in fetching payment by authority", operationName, err)
		}
		if payment == nil {
			return 0, paymentError.Payment_NotFound
		}
		resp, err := gateway.VerifyTransaction(dtos.VerifyTransactionDto{ID: payment.Authority, Amount: payment.Amount})
		if err != nil {
			return 0, types.NewServerError("Error in verifying transaction from server", operationName, err)
		}
		var transactionID uint
		if resp.IsSuccess {
			transaction, err := svc.createPaymentTransaction(tx, payment, dto.Gateway)
			if err != nil {
				return 0, err
			}
			transactionID = transaction.ID
			if err := svc.updateSuccessPayment(tx, payment.ID, transaction.ID, resp.RefCode); err != nil {
				return 0, err
			}
			if err := svc.updateSuccessOrder(tx, payment.OrderID); err != nil {
				return 0, err
			}
			if err := svc.createCourseParticipates(tx, payment.UserID, payment.OrderID); err != nil {
				return 0, err
			}
		} else {
			if err := svc.updateFailedOrder(tx, payment.OrderID); err != nil {
				return 0, err
			}
			if err := svc.updateFailedPayment(tx, payment.OrderID); err != nil {
				return 0, err
			}
		}
		return transactionID, nil
	})
	if err != nil {
		return err
	}
	log.Printf("New Transaction Generated : %v \n", transactionID)
	return nil
}

func (svc paymentService) selectGateway(paymentGateway entities.PaymentGateway) contracts.PaymentGateway {
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

func (svc paymentService) getMerchantID(paymentGateway entities.PaymentGateway) string {
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

func (svc paymentService) getCurrency(paymentGateway entities.PaymentGateway) string {
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

func (svc paymentService) FetchPageable(page, pageSize int) ([]*entities.Payment, int, error) {
	const operationName = "paymentService.FetchPageable"
	payments, count, err := svc.unitOfWork.PaymentRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset: &page,
			Limit:  &pageSize,
		},
	)
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching payments pages", operationName, err)
	}
	return payments, count, nil
}

func (svc paymentService) createPaymentTransaction(tx db.UnitOfWorkTx, payment *entities.Payment, gateway entities.PaymentGateway,
) (*entities.Transaction, error) {
	const operationName = "paymentService.createPaymentTransaction"
	transaction := &entities.Transaction{
		Amount:   payment.Amount,
		Type:     entities.TransactionType_Deposit,
		User:     payment.User.FullName(),
		Phone:    payment.User.Phone,
		Tag:      entities.TransactionTag_Sell,
		Currency: svc.getCurrency(gateway),
	}
	if err := tx.TransactionRepo().Create(transaction); err != nil {
		return nil, types.NewServerError("Error in creating transaction based on the payment", operationName, err)
	}
	return transaction, nil
}

func (svc paymentService) updateSuccessPayment(tx db.UnitOfWorkTx, paymentID, transactionID uint, refCode string) error {
	const operationName = "paymentService.updateSuccessPayment"
	payment := &entities.Payment{
		Model:           gorm.Model{ID: paymentID},
		TransactionID:   &transactionID,
		Status:          entities.PaymentStatus_Success,
		RefID:           refCode,
		StatusChangedAt: utils.Now(),
	}
	if err := tx.PaymentRepo().Update(payment); err != nil {
		return types.NewServerError("Error in updating successful payment", operationName, err)
	}
	return nil
}

func (svc paymentService) updateSuccessOrder(tx db.UnitOfWorkTx, orderID uint) error {
	const operationName = "paymentService.updateSuccessOrder"
	order := &entities.Order{
		Model:           gorm.Model{ID: orderID},
		Status:          entities.OrderStatus_Success,
		StatusChangedAt: utils.Now(),
	}
	if err := tx.OrderRepo().Update(order); err != nil {
		return types.NewServerError("Error in updating successful order", operationName, err)
	}
	return nil
}

func (svc paymentService) updateFailedOrder(tx db.UnitOfWorkTx, orderID uint) error {
	const operationName = "paymentService.updateFailedOrder"
	order := &entities.Order{
		Model:           gorm.Model{ID: orderID},
		Status:          entities.OrderStatus_Failed,
		StatusChangedAt: utils.Now(),
	}
	if err := tx.OrderRepo().Update(order); err != nil {
		return types.NewServerError("Error in updating failed order", operationName, err)
	}
	return nil
}

func (svc paymentService) updateFailedPayment(tx db.UnitOfWorkTx, paymentID uint) error {
	const operationName = "paymentService.updateFailedPayment"
	payment := &entities.Payment{
		Model:           gorm.Model{ID: paymentID},
		Status:          entities.PaymentStatus_Failure,
		StatusChangedAt: utils.Now(),
	}
	if err := tx.PaymentRepo().Update(payment); err != nil {
		return types.NewServerError("Error in updating failed payment", operationName, err)
	}
	return nil
}

func (svc paymentService) createCourseParticipates(tx db.UnitOfWorkTx, userID, orderID uint) error {
	const operationName = "paymentService.createCourseParticipates"
	order, err := tx.OrderRepo().GetByID(orderID, []string{"Items", "Items.Course"})
	if err != nil {
		return types.NewServerError("Error in getting order", operationName, err)
	}
	for _, item := range order.Items {
		courseParticipant := &entities.CourseParticipant{
			CourseID:  item.CourseID,
			TeacherID: *item.Course.TeacherID,
			StudentID: userID,
		}
		if err := tx.CourseParticipantRepo().Create(courseParticipant); err != nil {
			return types.NewServerError("Error in creating course participates", operationName, err)
		}
	}
	return nil
}
