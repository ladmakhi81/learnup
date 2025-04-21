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
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewCartService(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *CartServiceImpl {
	return &CartServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
	}
}

func (svc CartServiceImpl) Create(dto cartDtoReq.CreateCartReq) (*entities.Cart, error) {
	isCartExist, cartExistErr := svc.repo.CartRepo.Exist(map[string]any{
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
	course, courseErr := svc.repo.CourseRepo.GetByID(dto.CourseID)
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
	if err := svc.repo.CartRepo.Create(cart); err != nil {
		return nil, types.NewServerError(
			"Error in creating cart items",
			"CartServiceImpl.Create",
			err,
		)
	}
	return cart, nil
}

func (svc CartServiceImpl) DeleteByID(userID, id uint) error {
	cart, cartErr := svc.repo.CartRepo.GetByID(id)
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
	deleteErr := svc.repo.CartRepo.Delete(cart)
	if deleteErr != nil {
		return types.NewServerError(
			"Error in deleting cart by id",
			"CartServiceImpl.DeleteByID",
			deleteErr,
		)
	}
	return nil
}

func (svc CartServiceImpl) FetchAllByUserID(userID uint) ([]*entities.Cart, error) {
	user, userErr := svc.repo.UserRepo.GetByID(userID)
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
	carts, cartsErr := svc.repo.CartRepo.GetAll(repositories.GetAllOptions{
		Conditions: map[string]any{
			"user_id": userID,
		},
	})
	if cartsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching all carts by user id",
			"CartServiceImpl.FetchAllByUserID",
			cartsErr,
		)
	}
	return carts, nil
}
