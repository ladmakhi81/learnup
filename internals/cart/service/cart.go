package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	cartDtoReq "github.com/ladmakhi81/learnup/internals/cart/dto/req"
	cartRepo "github.com/ladmakhi81/learnup/internals/cart/repo"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type CartService interface {
	Create(dto cartDtoReq.CreateCartReq) (*entities.Cart, error)
	FetchByID(id uint) (*entities.Cart, error)
	DeleteByID(userID, id uint) error
	DeleteAllByUserID(userID uint) error
	FetchAllByUserID(userID uint) ([]*entities.Cart, error)
	FetchByUserAndCourse(userID uint, courseID uint) (*entities.Cart, error)
	FetchByCartIDs(ids []uint) ([]*entities.Cart, error)
}

type CartServiceImpl struct {
	cartRepo       cartRepo.CartRepo
	translationSvc contracts.Translator
	userSvc        userService.UserSvc
	courseSvc      courseService.CourseService
}

func NewCartService(
	cartRepo cartRepo.CartRepo,
	translationSvc contracts.Translator,
	userSvc userService.UserSvc,
	courseSvc courseService.CourseService,
) *CartServiceImpl {
	return &CartServiceImpl{
		cartRepo:       cartRepo,
		translationSvc: translationSvc,
		userSvc:        userSvc,
		courseSvc:      courseSvc,
	}
}

func (svc CartServiceImpl) Create(dto cartDtoReq.CreateCartReq) (*entities.Cart, error) {
	cartExist, cartExistErr := svc.FetchByUserAndCourse(dto.UserID, dto.CourseID)
	if cartExistErr != nil {
		return nil, cartExistErr
	}
	if cartExist != nil {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("cart.errors.exist_before"),
		)
	}
	course, courseErr := svc.courseSvc.FindById(dto.CourseID)
	if courseErr != nil {
		return nil, courseErr
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
	if err := svc.cartRepo.Create(cart); err != nil {
		return nil, types.NewServerError(
			"Error in creating cart items",
			"CartServiceImpl.Create",
			err,
		)
	}
	return cart, nil
}

func (svc CartServiceImpl) FetchByID(id uint) (*entities.Cart, error) {
	cart, cartErr := svc.cartRepo.FetchByID(id)
	if cartErr != nil {
		return nil, types.NewServerError(
			"Error in fetching cart by id",
			"CartServiceImpl.FetchByID",
			cartErr,
		)
	}
	return cart, nil
}

func (svc CartServiceImpl) DeleteByID(userID, id uint) error {
	cart, cartErr := svc.FetchByID(id)
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
	deleteErr := svc.cartRepo.DeleteByID(id)
	if deleteErr != nil {
		return types.NewServerError(
			"Error in deleting cart by id",
			"CartServiceImpl.DeleteByID",
			deleteErr,
		)
	}
	return nil
}

func (svc CartServiceImpl) DeleteAllByUserID(userID uint) error {
	user, userErr := svc.userSvc.FindById(userID)
	if userErr != nil {
		return userErr
	}
	if user == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	deleteErr := svc.cartRepo.DeleteAllByUserID(user.ID)
	if deleteErr != nil {
		return types.NewServerError(
			"Error in deleting carts by user id",
			"CartServiceImpl.DeleteAllByUserID",
			deleteErr,
		)
	}
	return nil
}

func (svc CartServiceImpl) FetchAllByUserID(userID uint) ([]*entities.Cart, error) {
	user, userErr := svc.userSvc.FindById(userID)
	if userErr != nil {
		return nil, userErr
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	carts, cartsErr := svc.cartRepo.FetchAllByUserID(user.ID)
	if cartsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching all carts by user id",
			"CartServiceImpl.FetchAllByUserID",
			cartsErr,
		)
	}
	return carts, nil
}

func (svc CartServiceImpl) FetchByUserAndCourse(userID uint, courseID uint) (*entities.Cart, error) {
	cart, cartErr := svc.cartRepo.FetchByUserAndCourse(userID, courseID)
	if cartErr != nil {
		return nil, types.NewServerError(
			"Error in fetching cart by user and course",
			"CartServiceImpl.FetchByUserAndCourse",
			cartErr,
		)
	}

	return cart, nil
}

func (svc CartServiceImpl) FetchByCartIDs(ids []uint) ([]*entities.Cart, error) {
	carts, cartsErr := svc.cartRepo.FetchByCartIDs(ids)
	if cartsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching carts by ids",
			"CartServiceImpl.FetchByCartIDS",
			cartsErr,
		)
	}
	return carts, nil
}
