package service

import (
	cartError "github.com/ladmakhi81/learnup/internals/cart/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	orderDtoReq "github.com/ladmakhi81/learnup/internals/order/dto/req"
	paymentDtoReq "github.com/ladmakhi81/learnup/internals/payment/dto/req"
	paymentService "github.com/ladmakhi81/learnup/internals/payment/service"
	userError "github.com/ladmakhi81/learnup/internals/user/error"
	"github.com/ladmakhi81/learnup/types"
)

type OrderService interface {
	Create(dto orderDtoReq.CreateOrderReq) (string, error)
	FetchPaginated(page, pageSize int) ([]*entities.Order, int, error)
	FetchDetailById(id uint) (*entities.Order, error)
}

type orderService struct {
	unitOfWork db.UnitOfWork
	paymentSvc paymentService.PaymentService
}

func NewOrderService(unitOfWork db.UnitOfWork, paymentSvc paymentService.PaymentService) OrderService {
	return &orderService{unitOfWork: unitOfWork, paymentSvc: paymentSvc}
}

func (svc orderService) Create(dto orderDtoReq.CreateOrderReq) (string, error) {
	const operationName = "orderService.Create"
	return db.WithTx(svc.unitOfWork, func(tx db.UnitOfWorkTx) (string, error) {
		user, err := tx.UserRepo().GetByID(dto.UserID, nil)
		if err != nil {
			return "", types.NewServerError("Error in fetching user by id", operationName, err)
		}
		if user == nil {
			return "", userError.User_NotFound
		}
		carts, err := tx.CartRepo().GetAll(repositories.GetAllOptions{
			Conditions: map[string]any{
				"id": dto.Carts,
			},
			Relations: []string{"Course"},
		})
		if err != nil {
			return "", types.NewServerError("Error in fetching all carts based on carts ids", operationName, err)
		}
		if len(carts) != len(dto.Carts) || len(carts) == 0 {
			return "", cartError.Cart_ListNotMatch
		}
		order := entities.NewOrder(user.ID)
		if err := tx.OrderRepo().Create(order); err != nil {
			return "", types.NewServerError("Error in creating order", operationName, err)
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
		if err := tx.OrderItemRepo().BatchInsert(orderItems); err != nil {
			return "", types.NewServerError("Error in batch insert order items", operationName, err)
		}
		order.TotalPrice = totalAmount
		order.FinalPrice = totalAmount
		if err := tx.OrderRepo().Update(order); err != nil {
			return "", types.NewServerError("Error in updating order", operationName, err)
		}
		payment, err := svc.paymentSvc.Create(
			tx,
			paymentDtoReq.CreatePaymentReq{
				Gateway: dto.Gateway,
				UserID:  user.ID,
				OrderID: order.ID,
				Amount:  totalAmount,
			},
		)
		if err != nil {
			return "", err
		}
		if err := tx.CartRepo().BatchDelete(carts); err != nil {
			return "", types.NewServerError("Error in deleting carts as batching", operationName, err)
		}
		return payment.PayLink, nil
	})
}

func (svc orderService) FetchPaginated(page, pageSize int) ([]*entities.Order, int, error) {
	const operationName = "orderService.FetchPaginated"
	orders, count, err := svc.unitOfWork.OrderRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Relations: []string{
			"User",
		},
	})
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching paginated list of orders", operationName, err)
	}
	return orders, count, nil
}

func (svc orderService) FetchDetailById(id uint) (*entities.Order, error) {
	const operationName = "orderService.FetchDetailById"
	order, err := svc.unitOfWork.OrderRepo().GetByID(id, []string{"User", "Items", "Items.Course"})
	if err != nil {
		return nil, types.NewServerError("Error in fetching detail by id", operationName, err)
	}
	return order, nil
}
