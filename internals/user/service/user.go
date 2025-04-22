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
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	isPhoneExistBefore, isPhoneExistBeforeErr := tx.UserRepo().Exist(map[string]any{
		"phone_number": dto.Phone,
	})
	if isPhoneExistBeforeErr != nil {
		return nil, types.NewServerError(
			"Error in checking phone number exist",
			"UserSvcImpl.CreateBasic",
			isPhoneExistBeforeErr,
		)
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
	if err := tx.UserRepo().Create(user); err != nil {
		return nil, types.NewServerError(
			"Create Basic User Throw Error",
			"UserSvcImpl.CreateBasic",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}
