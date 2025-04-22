package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/user/dto/req"
	userError "github.com/ladmakhi81/learnup/internals/user/error"
	"github.com/ladmakhi81/learnup/types"
	"golang.org/x/crypto/bcrypt"
)

type UserSvc interface {
	CreateBasic(dto dtoreq.CreateBasicUserReq) (*entities.User, error)
}

type userService struct {
	unitOfWork db.UnitOfWork
}

func NewUserSvc(unitOfWork db.UnitOfWork) UserSvc {
	return &userService{unitOfWork: unitOfWork}
}

func (svc userService) CreateBasic(dto dtoreq.CreateBasicUserReq) (*entities.User, error) {
	const operationName = "userService.CreateBasic"
	isPhoneExistBefore, err := svc.unitOfWork.UserRepo().Exist(map[string]any{"phone_number": dto.Phone})
	if err != nil {
		return nil, types.NewServerError("Error in checking phone number exist", operationName, err)
	}
	if isPhoneExistBefore {
		return nil, userError.User_PhoneDuplicated
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, types.NewServerError("Generating Password Throw Error", operationName, err)
	}
	user := &entities.User{
		Phone:     dto.Phone,
		Password:  string(hashedPassword),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
	if err := svc.unitOfWork.UserRepo().Create(user); err != nil {
		return nil, types.NewServerError("Create Basic User Throw Error", operationName, err)
	}
	return user, nil
}
