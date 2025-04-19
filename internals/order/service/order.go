package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	cartService "github.com/ladmakhi81/learnup/internals/cart/service"
	orderDtoReq "github.com/ladmakhi81/learnup/internals/order/dto/req"
	orderRepository "github.com/ladmakhi81/learnup/internals/order/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type OrderService interface {
	Create(dto orderDtoReq.CreateOrderReq) error
	FetchPaginated(page, pageSize int) ([]*entities.Order, error)
	FetchCount() (int, error)
	FetchDetailById(id uint) (*entities.Order, error)
}

type OrderServiceImpl struct {
	orderRepo      orderRepository.OrderRepo
	orderItemRepo  orderRepository.OrderItemRepo
	userSvc        userService.UserSvc
	cartSvc        cartService.CartService
	translationSvc contracts.Translator
}

func NewOrderService(
	orderRepo orderRepository.OrderRepo,
	orderItemRepo orderRepository.OrderItemRepo,
	userSvc userService.UserSvc,
	cartSvc cartService.CartService,
	translationSvc contracts.Translator,
) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo:      orderRepo,
		userSvc:        userSvc,
		cartSvc:        cartSvc,
		translationSvc: translationSvc,
		orderItemRepo:  orderItemRepo,
	}
}

func (svc OrderServiceImpl) Create(dto orderDtoReq.CreateOrderReq) error {
	user, userErr := svc.userSvc.FindById(dto.UserID)
	if userErr != nil {
		return userErr
	}
	if user == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	carts, cartsErr := svc.cartSvc.FetchByCartIDs(dto.Carts)
	if cartsErr != nil {
		return cartsErr
	}
	if len(carts) != len(dto.Carts) || len(carts) == 0 {
		return types.NewNotFoundError(
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
	if err := svc.orderRepo.Create(order); err != nil {
		return types.NewServerError(
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
	if err := svc.orderItemRepo.CreateBatch(orderItems); err != nil {
		return types.NewServerError(
			"Error in batch insert order items",
			"OrderItemRepo.CreateBatch",
			err,
		)
	}
	order.TotalPrice = totalAmount
	if err := svc.orderRepo.Update(order); err != nil {
		return types.NewServerError(
			"Error in updating order",
			"OrderServiceImpl.Create",
			err,
		)
	}
	return nil
}

func (svc OrderServiceImpl) FetchPaginated(page, pageSize int) ([]*entities.Order, error) {
	orders, ordersErr := svc.orderRepo.FetchPaginated(page, pageSize)
	if ordersErr != nil {
		return nil, types.NewServerError(
			"Error in fetching paginated list of orders",
			"OrderServiceImpl.FetchPaginated",
			ordersErr,
		)
	}
	return orders, nil
}

func (svc OrderServiceImpl) FetchCount() (int, error) {
	count, countErr := svc.orderRepo.FetchCount()
	if countErr != nil {
		return 0, types.NewServerError(
			"Error in fetching count",
			"OrderServiceImpl.FetchCount",
			countErr,
		)
	}
	return count, nil
}

func (svc OrderServiceImpl) FetchDetailById(id uint) (*entities.Order, error) {
	order, orderErr := svc.orderRepo.FetchDetailById(id)
	if orderErr != nil {
		return nil, types.NewServerError(
			"Error in fetching detail by id",
			"OrderServiceImpl.FetchDetailById",
			orderErr,
		)
	}
	return order, nil
}
