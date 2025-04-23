package repositories

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
)

type TransactionRepo interface {
	Repository[entities.Transaction]
}

type TransactionRepoImpl struct {
	RepositoryImpl[entities.Transaction]
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepoImpl {
	return &TransactionRepoImpl{
		RepositoryImpl[entities.Transaction]{
			db: db,
		},
	}
}
