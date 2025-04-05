package service

import (
	"github.com/ladmakhi81/learnup/internals/auth/constant"
	dtoreq "github.com/ladmakhi81/learnup/internals/auth/dto/req"
	"github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/cache"
	"github.com/ladmakhi81/learnup/pkg/token"
	"github.com/ladmakhi81/learnup/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req dtoreq.LoginReq) (string, error)
}

type AuthServiceImpl struct {
	userSvc  service.UserSvc
	cacheSvc cache.Cache
	tokenSvc token.Token
}

func NewAuthServiceImpl(
	userSvc service.UserSvc,
	cacheSvc cache.Cache,
	tokenSvc token.Token,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		userSvc:  userSvc,
		cacheSvc: cacheSvc,
		tokenSvc: tokenSvc,
	}
}

func (svc *AuthServiceImpl) Login(dto dtoreq.LoginReq) (string, error) {
	user, userErr := svc.userSvc.FindByPhone(dto.Phone)
	if userErr != nil {
		return "", userErr
	}
	if user == nil {
		return "", types.NewNotFoundError(constant.AuthInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return "", types.NewNotFoundError(constant.AuthInvalidCredentials)
	}
	accessToken, accessTokenErr := svc.tokenSvc.GenerateToken(user.ID)
	if accessTokenErr != nil {
		return "", types.NewServerError(
			"Error in generating access token",
			"AuthServiceImpl.Login",
			accessTokenErr,
		)
	}
	if err := svc.cacheSvc.SetVal(constant.LoginCacheKey, accessToken); err != nil {
		return "", types.NewServerError(
			"Error in updating cache redis",
			"AuthServiceImpl.Login",
			err,
		)
	}
	return accessToken, nil
}
