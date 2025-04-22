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
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	isCartExist, cartExistErr := tx.CartRepo().Exist(map[string]any{
		"course_id": dto.CourseID,
		"user_id":   dto.UserID,
	})
	if cartExistErr != nil {
		return nil, types.NewServerError(
			"Error in checking is cart exist or not",
			"CartServiceImpl.Create",
			cartExistErr,
		)
	}
	if isCartExist {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("cart.errors.exist_before"),
		)
	}
	course, courseErr := tx.CourseRepo().GetByID(dto.CourseID, nil)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in fetching course by id",
			"CartServiceImpl.Create",
			courseErr,
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
	if err := tx.CartRepo().Create(cart); err != nil {
		return nil, types.NewServerError(
			"Error in creating cart items",
			"CartServiceImpl.Create",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return cart, nil
}

func (svc CartServiceImpl) DeleteByID(userID, id uint) error {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return txErr
	}
	cart, cartErr := tx.CartRepo().GetByID(id, nil)
	if cartErr != nil {
		return cartErr
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
	deleteErr := tx.CartRepo().Delete(cart)
	if deleteErr != nil {
		return types.NewServerError(
			"Error in deleting cart by id",
			"CartServiceImpl.DeleteByID",
			deleteErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (svc CartServiceImpl) FetchAllByUserID(userID uint) ([]*entities.Cart, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	user, userErr := tx.UserRepo().GetByID(userID, nil)
	if userErr != nil {
		return nil, types.NewServerError(
			"Error in fetching carts by user id",
			"CartServiceImpl.FetchAllByUserID",
			userErr,
		)
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	carts, cartsErr := tx.CartRepo().GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"user_id": userID,
		},
		Relations: []string{"Course"},
	})
	if cartsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching all carts by user id",
			"CartServiceImpl.FetchAllByUserID",
			cartsErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return carts, nil
}
