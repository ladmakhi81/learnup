package repositories

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"gorm.io/gorm"
)

type PaymentRepo interface {
	Repository[entities.Payment]
}

type PaymentRepoImpl struct {
	RepositoryImpl[entities.Payment]
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepoImpl {
	return &PaymentRepoImpl{
		RepositoryImpl[entities.Payment]{
			db: db,
		},
	}
}
