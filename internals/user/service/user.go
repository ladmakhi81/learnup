package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/user/dto/req"
	"github.com/ladmakhi81/learnup/internals/user/repo"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"golang.org/x/crypto/bcrypt"
)

type UserSvc interface {
	CreateBasic(dto dtoreq.CreateBasicUserReq) (*entities.User, error)
	IsPhoneDuplicated(phone string) (bool, error)
	FindByPhone(phone string) (*entities.User, error)
	FindById(id uint) (*entities.User, error)
}

type UserSvcImpl struct {
	userRepo       repo.UserRepo
	translationSvc contracts.Translator
}

func NewUserSvcImpl(
	userRepo repo.UserRepo,
	translationSvc contracts.Translator,
) *UserSvcImpl {
	return &UserSvcImpl{
		userRepo:       userRepo,
		translationSvc: translationSvc,
	}
}

func (svc UserSvcImpl) CreateBasic(dto dtoreq.CreateBasicUserReq) (*entities.User, error) {
	isPhoneExistBefore, isPhoneExistBeforeErr := svc.IsPhoneDuplicated(dto.Phone)
	if isPhoneExistBeforeErr != nil {
		return nil, isPhoneExistBeforeErr
	}
	if isPhoneExistBefore {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("user.errors.phone_duplicate"),
		)
	}
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		return nil, types.NewServerError(
			"Generating Password Throw Error",
			"UserSvcImpl.CreateBasic",
			hashedPasswordErr,
		)
	}
	user := &entities.User{
		Phone:     dto.Phone,
		Password:  string(hashedPassword),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
	if err := svc.userRepo.CreateBasic(user); err != nil {
		return nil, types.NewServerError(
			"Create Basic User Throw Error",
			"UserSvcImpl.CreateBasic",
			err,
		)
	}
	return user, nil
}

func (svc UserSvcImpl) IsPhoneDuplicated(phone string) (bool, error) {
	user, userErr := svc.FindByPhone(phone)
	if userErr != nil {
		return false, userErr
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}

func (svc UserSvcImpl) FindByPhone(phone string) (*entities.User, error) {
	user, userErr := svc.userRepo.FetchByPhone(phone)
	if userErr != nil {
		return nil, types.NewServerError(
			"Find User By Phone Throw Error",
			"UserSvcImpl.FetchByPhone",
			userErr,
		)
	}
	return user, nil
}

func (svc UserSvcImpl) FindById(id uint) (*entities.User, error) {
	user, userErr := svc.userRepo.FetchById(id)
	if userErr != nil {
		return nil, types.NewServerError(
			"Find User By Id Throw Error",
			"UserServiceImpl.FetchById",
			userErr,
		)
	}
	return user, nil
}
