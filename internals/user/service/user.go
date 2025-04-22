package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/user/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"golang.org/x/crypto/bcrypt"
)

type UserSvc interface {
	CreateBasic(dto dtoreq.CreateBasicUserReq) (*entities.User, error)
}

type UserSvcImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewUserSvcImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *UserSvcImpl {
	return &UserSvcImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc UserSvcImpl) CreateBasic(dto dtoreq.CreateBasicUserReq) (*entities.User, error) {
	const operationName = "UserSvcImpl.CreateBasic"
	isPhoneExistBefore, err := svc.unitOfWork.UserRepo().Exist(map[string]any{
		"phone_number": dto.Phone,
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in checking phone number exist",
			operationName,
			err,
		)
	}
	if isPhoneExistBefore {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("user.errors.phone_duplicate"),
		)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, types.NewServerError(
			"Generating Password Throw Error",
			operationName,
			err,
		)
	}
	user := &entities.User{
		Phone:     dto.Phone,
		Password:  string(hashedPassword),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
	if err := svc.unitOfWork.UserRepo().Create(user); err != nil {
		return nil, types.NewServerError(
			"Create Basic User Throw Error",
			operationName,
			err,
		)
	}
	return user, nil
}
