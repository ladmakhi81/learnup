package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/user/entity"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateBasic(user *entity.User) error
	FindByPhone(phone string) (*entity.User, error)
}

type UserRepoImpl struct {
	db *db.Database
}

func NewUserRepoImpl(db *db.Database) *UserRepoImpl {
	return &UserRepoImpl{
		db: db,
	}
}

func (svc UserRepoImpl) CreateBasic(user *entity.User) error {
	tx := svc.db.Core.Create(user)
	return tx.Error
}

func (svc UserRepoImpl) FindByPhone(phone string) (*entity.User, error) {
	user := new(entity.User)
	tx := svc.db.Core.Where("phone = ?", phone).First(user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return user, nil
}
