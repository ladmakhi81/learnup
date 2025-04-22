package service

import (
	"github.com/ladmakhi81/learnup/internals/auth/constant"
	dtoreq "github.com/ladmakhi81/learnup/internals/auth/dto/req"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req dtoreq.LoginReq) (string, error)
}

type AuthServiceImpl struct {
	cacheSvc       contracts.Cache
	tokenSvc       contracts.Token
	translationSvc contracts.Translator
	unitOfWork     db.UnitOfWork
}

func NewAuthServiceImpl(
	cacheSvc contracts.Cache,
	tokenSvc contracts.Token,
	translationSvc contracts.Translator,
	unitOfWork db.UnitOfWork,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		cacheSvc:       cacheSvc,
		tokenSvc:       tokenSvc,
		translationSvc: translationSvc,
		unitOfWork:     unitOfWork,
	}
}

func (svc AuthServiceImpl) Login(dto dtoreq.LoginReq) (string, error) {
	const operationName = "AuthServiceImpl.Login"
	user, err := svc.unitOfWork.UserRepo().GetOne(map[string]any{
		"phone_number": dto.Phone,
	}, nil)
	if err != nil {
		return "", types.NewServerError(
			"Error in fetching user by phone number",
			operationName,
			err,
		)
	}
	if user == nil {
		return "", types.NewNotFoundError(
			svc.translationSvc.Translate("auth.errors.invalid_credentials"),
		)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return "", types.NewNotFoundError(
			svc.translationSvc.Translate("auth.errors.invalid_credentials"),
		)
	}
	accessToken, err := svc.tokenSvc.GenerateToken(user.ID)
	if err != nil {
		return "", types.NewServerError(
			"Error in generating access token",
			operationName,
			err,
		)
	}
	if err := svc.cacheSvc.SetVal(constant.LoginCacheKey, accessToken); err != nil {
		return "", types.NewServerError(
			"Error in updating cache redis",
			operationName,
			err,
		)
	}
	return accessToken, nil
}
