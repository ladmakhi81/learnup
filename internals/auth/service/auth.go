package service

import (
	"github.com/ladmakhi81/learnup/internals/auth/constant"
	dtoreq "github.com/ladmakhi81/learnup/internals/auth/dto/req"
	authError "github.com/ladmakhi81/learnup/internals/auth/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type AuthService interface {
	Login(req dtoreq.LoginReq) (string, error)
}

type authService struct {
	cacheSvc   contracts.Cache
	tokenSvc   contracts.Token
	unitOfWork db.UnitOfWork
}

func NewAuthSvc(
	cacheSvc contracts.Cache,
	tokenSvc contracts.Token,
	unitOfWork db.UnitOfWork,
) AuthService {
	return &authService{
		cacheSvc:   cacheSvc,
		tokenSvc:   tokenSvc,
		unitOfWork: unitOfWork,
	}
}

func (svc authService) Login(dto dtoreq.LoginReq) (string, error) {
	const operationName = "authService.Login"
	user, err := svc.unitOfWork.UserRepo().GetOne(map[string]any{"phone_number": dto.Phone}, nil)
	if err != nil {
		return "", types.NewServerError("Error in fetching user by phone number", operationName, err)
	}
	if user == nil || !user.IsPasswordMatch(dto.Password) {
		return "", authError.Auth_InvalidCredentials
	}
	accessToken, err := svc.tokenSvc.GenerateToken(user.ID)
	if err != nil {
		return "", types.NewServerError("Error in generating access token", operationName, err)
	}
	if err := svc.cacheSvc.SetVal(constant.LoginCacheKey, accessToken); err != nil {
		return "", types.NewServerError("Error in updating cache redis", operationName, err)
	}
	return accessToken, nil
}
