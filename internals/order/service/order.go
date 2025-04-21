package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	orderDtoReq "github.com/ladmakhi81/learnup/internals/order/dto/req"
	paymentDtoReq "github.com/ladmakhi81/learnup/internals/payment/dto/req"
	paymentService "github.com/ladmakhi81/learnup/internals/payment/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type OrderService interface {
	Create(dto orderDtoReq.CreateOrderReq) (string, error)
	FetchPaginated(page, pageSize int) ([]*entities.Order, int, error)
	FetchDetailById(id uint) (*entities.Order, error)
}

type OrderServiceImpl struct {
	repo           *db.Repositories
	translationSvc contracts.Translator
	paymentSvc     paymentService.PaymentService
}

func NewOrderService(
	repo *db.Repositories,
	translationSvc contracts.Translator,
	paymentSvc paymentService.PaymentService,
) *OrderServiceImpl {
	return &OrderServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
		paymentSvc:     paymentSvc,
	}
}

func (svc OrderServiceImpl) Create(dto orderDtoReq.CreateOrderReq) (string, error) {
	user, userErr := svc.repo.UserRepo.GetByID(dto.UserID)
	if userErr != nil {
		return "", types.NewServerError(
			"Error in fetching user by id",
			"OrderServiceImpl.Create",
			userErr,
		)
	}
	if user == nil {
		return "", types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	carts, cartsErr := svc.repo.CartRepo.GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"id": dto.Carts,
		},
	})
	if cartsErr != nil {
		return "", types.NewServerError(
			"Error in fetching all carts based on carts ids",
			"OrderServiceImpl.Create",
			cartsErr,
		)
	}
	if len(carts) != len(dto.Carts) || len(carts) == 0 {
		return "", types.NewNotFoundError(
			svc.translationSvc.Translate("cart.errors.list_not_match"),
		)
	}
	order := &entities.Order{
		UserID:        dto.UserID,
		TotalPrice:    0,
		DiscountPrice: 0,
		FinalPrice:    0,
		Status:        entities.OrderStatus_Pending,
	}
	if err := svc.repo.OrderRepo.Create(order); err != nil {
		return "", types.NewServerError(
			"Error in creating order",
			"OrderServiceImpl.Create",
			err,
		)
	}
	var totalAmount float64
	orderItems := make([]*entities.OrderItem, len(carts))
	for index, cart := range carts {
		orderItems[index] = &entities.OrderItem{
			UserID:   user.ID,
			CourseID: cart.CourseID,
			OrderID:  order.ID,
			Amount:   cart.Course.Price,
		}
		totalAmount += cart.Course.Price
	}
	if err := svc.repo.OrderItemRepo.BatchInsert(orderItems); err != nil {
		return "", types.NewServerError(
			"Error in batch insert order items",
			"OrderServiceImpl.Create",
			err,
		)
	}
	order.TotalPrice = totalAmount
	order.FinalPrice = totalAmount
	if err := svc.repo.OrderRepo.Update(order); err != nil {
		return "", types.NewServerError(
			"Error in updating order",
			"OrderServiceImpl.Create",
			err,
		)
	}
	payment, paymentErr := svc.paymentSvc.Create(paymentDtoReq.CreatePaymentReq{
		Gateway: dto.Gateway,
		UserID:  user.ID,
		OrderID: order.ID,
		Amount:  totalAmount,
	})
	if paymentErr != nil {
		return "", paymentErr
	}
	if err := svc.repo.CartRepo.BatchDelete(carts); err != nil {
		return "", types.NewServerError(
			"Error in deleting carts as batching",
			"OrderServiceImpl.Create",
			err,
		)
	}
	return payment.PayLink, nil
}

func (svc OrderServiceImpl) FetchPaginated(page, pageSize int) ([]*entities.Order, int, error) {
	orders, count, ordersErr := svc.repo.OrderRepo.GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Relations: []string{
			"User",
		},
	})
	if ordersErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching paginated list of orders",
			"OrderServiceImpl.FetchPaginated",
			ordersErr,
		)
	}
	return orders, count, nil
}

func (svc OrderServiceImpl) FetchDetailById(id uint) (*entities.Order, error) {
	//TODO: add relation to fetch detail
	//Preload("User").
	//Preload("Items").
	//Preload("Items.Course").
	order, orderErr := svc.repo.OrderRepo.GetByID(id)
	if orderErr != nil {
		return nil, types.NewServerError(
			"Error in fetching detail by id",
			"OrderServiceImpl.FetchDetailById",
			orderErr,
		)
	}
	return order, nil
}
