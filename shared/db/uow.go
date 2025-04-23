package db

import (
	"github.com/ladmakhi81/learnup/shared/types"
	"gorm.io/gorm"
)

type UnitOfWork interface {
	Begin() (UnitOfWorkTx, error)
	Repo
}

type UnitOfWorkTx interface {
	Repo
	Commit() error
	Rollback() error
}

type UnitOfWorkImpl struct {
	db *gorm.DB
	*RepoProvider
}

type UnitOfWorkTxImpl struct {
	tx *gorm.DB
	*RepoProvider
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &UnitOfWorkImpl{
		db:           db,
		RepoProvider: NewRepoProvider(db),
	}
}

func (svc UnitOfWorkImpl) Begin() (UnitOfWorkTx, error) {
	tx := svc.db.Begin()
	if tx.Error != nil {
		return nil, types.NewServerError(
			"Error in begin transaction",
			"UnitOfWorkImpl.Begin",
			tx.Error,
		)
	}
	return NewUnitOfWorkTx(tx), nil
}

func NewUnitOfWorkTx(tx *gorm.DB) UnitOfWorkTx {
	return &UnitOfWorkTxImpl{
		tx:           tx,
		RepoProvider: NewRepoProvider(tx),
	}
}

func (svc UnitOfWorkTxImpl) Commit() error {
	if err := svc.tx.Commit().Error; err != nil {
		return types.NewServerError(
			"Error in commit changes",
			"UnitOfWorkTxImpl.Commit",
			err,
		)
	}
	return nil
}

func (svc UnitOfWorkTxImpl) Rollback() error {
	if err := svc.tx.Rollback().Error; err != nil {
		return types.NewServerError(
			"Error in rollback",
			"UnitOfWorkTxImpl.Rollback",
			err,
		)
	}
	return nil
}

func WithTx[T any](unitOfWork UnitOfWork, fn func(tx UnitOfWorkTx) (T, error)) (T, error) {
	var defaultRes T
	tx, err := unitOfWork.Begin()
	if err != nil {
		return defaultRes, err
	}
	resp, err := fn(tx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	if err != nil {
		tx.Rollback()
		return defaultRes, err
	}

	tx.Commit()
	return resp, nil
}
