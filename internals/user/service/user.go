package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/user/dto/req"
	"github.com/ladmakhi81/learnup/internals/user/entity"
	"github.com/ladmakhi81/learnup/internals/user/repo"
)

type UserSvc interface {
	CreateBasic(dto dtoreq.CreateUserReq) (*entity.User, error)
	IsPhoneDuplicated(phone string) (bool, error)
	FindByPhone(phone string) (*entity.User, error)
}

type UserSvcImpl struct {
	userRepo repo.UserRepo
}

func NewUserSvcImpl(userRepo repo.UserRepo) *UserSvcImpl {
	return &UserSvcImpl{
		userRepo: userRepo,
	}
}

func (svc UserSvcImpl) CreateBasic(dto dtoreq.CreateUserReq) (*entity.User, error) {
	return nil, nil
}

func (svc UserSvcImpl) IsPhoneDuplicated(phone string) (bool, error) {
	return false, nil
}

func (svc UserSvcImpl) FindByPhone(phone string) (*entity.User, error) {
	return nil, nil
}
