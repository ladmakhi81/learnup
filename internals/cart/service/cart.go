package service

import (
	cartDtoReq "github.com/ladmakhi81/learnup/internals/cart/dto/req"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type CartService interface {
	Create(dto cartDtoReq.CreateCartReq) (*entities.Cart, error)
	DeleteByID(userID, id uint) error
	FetchAllByUserID(userID uint) ([]*entities.Cart, error)
}

type CartServiceImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewCartService(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *CartServiceImpl {
	return &CartServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc CartServiceImpl) Create(dto cartDtoReq.CreateCartReq) (*entities.Cart, error) {
	const operationName = "CartServiceImpl.Create"
	isCartExist, err := svc.unitOfWork.CartRepo().Exist(map[string]any{
		"course_id": dto.CourseID,
		"user_id":   dto.UserID,
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in checking is cart exist or not",
			operationName,
			err,
		)
	}
	if isCartExist {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("cart.errors.exist_before"),
		)
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching course by id",
			operationName,
			err,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	cart := &entities.Cart{
		UserID:   dto.UserID,
		CourseID: dto.CourseID,
	}
	if err := svc.unitOfWork.CartRepo().Create(cart); err != nil {
		return nil, types.NewServerError(
			"Error in creating cart items",
			operationName,
			err,
		)
	}
	return cart, nil
}

func (svc CartServiceImpl) DeleteByID(userID, id uint) error {
	const operationName = "CartServiceImpl.DeleteByID"
	cart, err := svc.unitOfWork.CartRepo().GetByID(id, nil)
	if err != nil {
		return err
	}
	if cart == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("cart.errors.not_found"),
		)
	}
	if cart.UserID != userID {
		return types.NewForbiddenAccessError(
			svc.translationSvc.Translate("cart.errors.owner_delete"),
		)
	}
	if err := svc.unitOfWork.CartRepo().Delete(cart); err != nil {
		return types.NewServerError(
			"Error in deleting cart by id",
			operationName,
			err,
		)
	}
	return nil
}

func (svc CartServiceImpl) FetchAllByUserID(userID uint) ([]*entities.Cart, error) {
	const operationName = "CartServiceImpl.FetchAllByUserID"
	user, err := svc.unitOfWork.UserRepo().GetByID(userID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching carts by user id",
			operationName,
			err,
		)
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	carts, err := svc.unitOfWork.CartRepo().GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"user_id": userID,
		},
		Relations: []string{"Course"},
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching all carts by user id",
			operationName,
			err,
		)
	}
	return carts, nil
}
