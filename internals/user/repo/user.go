package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateBasic(user *entities.User) error
	FetchByPhone(phone string) (*entities.User, error)
	FetchById(id uint) (*entities.User, error)
}

type UserRepoImpl struct {
	db *db.Database
}

func NewUserRepoImpl(db *db.Database) *UserRepoImpl {
	return &UserRepoImpl{
		db: db,
	}
}

func (svc UserRepoImpl) CreateBasic(user *entities.User) error {
	tx := svc.db.Core.Create(user)
	return tx.Error
}

func (svc UserRepoImpl) FetchByPhone(phone string) (*entities.User, error) {
	user := new(entities.User)
	tx := svc.db.Core.Where("phone_number = ?", phone).First(user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return user, nil
}

func (svc UserRepoImpl) FetchById(id uint) (*entities.User, error) {
	user := &entities.User{}
	tx := svc.db.Core.Where("id = ?", id).First(user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return user, nil
}
