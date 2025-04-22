package service

import (
	cartDtoReq "github.com/ladmakhi81/learnup/internals/cart/dto/req"
	cartError "github.com/ladmakhi81/learnup/internals/cart/error"
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	userError "github.com/ladmakhi81/learnup/internals/user/error"
	"github.com/ladmakhi81/learnup/types"
)

type CartService interface {
	Create(dto cartDtoReq.CreateCartReq) (*entities.Cart, error)
	DeleteByID(userID, id uint) error
	FetchAllByUserID(userID uint) ([]*entities.Cart, error)
}

type cartService struct {
	unitOfWork db.UnitOfWork
}

func NewCartSvc(unitOfWork db.UnitOfWork) CartService {
	return &cartService{unitOfWork: unitOfWork}
}

func (svc cartService) Create(dto cartDtoReq.CreateCartReq) (*entities.Cart, error) {
	const operationName = "cartService.Create"
	isCartExist, err := svc.unitOfWork.CartRepo().Exist(map[string]any{"course_id": dto.CourseID, "user_id": dto.UserID})
	if err != nil {
		return nil, types.NewServerError("Error in checking cart exist", operationName, err)
	}
	if isCartExist {
		return nil, cartError.Cart_Duplicated
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching course by id", operationName, err)
	}
	if course == nil {
		return nil, courseError.Course_NotFound
	}
	cart := &entities.Cart{
		UserID:   dto.UserID,
		CourseID: dto.CourseID,
	}
	if err := svc.unitOfWork.CartRepo().Create(cart); err != nil {
		return nil, types.NewServerError("Error in creating cart items", operationName, err)
	}
	return cart, nil
}

func (svc cartService) DeleteByID(userID, id uint) error {
	const operationName = "cartService.DeleteByID"
	cart, err := svc.unitOfWork.CartRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError("Error in fetching cart by id", operationName, err)
	}
	if cart == nil {
		return cartError.Cart_NotFound
	}
	if cart.IsOwner(userID) {
		return cartError.Cart_ForbiddenAccess
	}
	if err := svc.unitOfWork.CartRepo().Delete(cart); err != nil {
		return types.NewServerError("Error in deleting cart by id", operationName, err)
	}
	return nil
}

func (svc cartService) FetchAllByUserID(userID uint) ([]*entities.Cart, error) {
	const operationName = "cartService.FetchAllByUserID"
	user, err := svc.unitOfWork.UserRepo().GetByID(userID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching carts by user id", operationName, err)
	}
	if user == nil {
		return nil, userError.User_NotFound
	}
	carts, err := svc.unitOfWork.CartRepo().GetAll(repositories.GetAllOptions{Conditions: map[string]any{"user_id": userID}, Relations: []string{"Course"}})
	if err != nil {
		return nil, types.NewServerError("Error in fetching all carts by user id", operationName, err)
	}
	return carts, nil
}
