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
	repo           *db.Repositories
}

func NewAuthServiceImpl(
	cacheSvc contracts.Cache,
	tokenSvc contracts.Token,
	translationSvc contracts.Translator,
	repo *db.Repositories,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		cacheSvc:       cacheSvc,
		tokenSvc:       tokenSvc,
		translationSvc: translationSvc,
		repo:           repo,
	}
}

func (svc AuthServiceImpl) Login(dto dtoreq.LoginReq) (string, error) {
	user, userErr := svc.repo.UserRepo.GetOne(map[string]any{
		"phone_number": dto.Phone,
	}, nil)
	if userErr != nil {
		return "", types.NewServerError(
			"Error in fetching user by phone number",
			"AuthServiceImpl.Login",
			userErr,
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
