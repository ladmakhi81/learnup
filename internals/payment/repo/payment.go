package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type PaymentRepo interface {
	Create(payment *entities.Payment) error
	FetchByAuthority(authority string) (*entities.Payment, error)
	Update(payment *entities.Payment) error
	FetchPageable(page, pageSize int) ([]*entities.Payment, error)
	FetchCount() (int, error)
}

type PaymentRepoImpl struct {
	dbClient *db.Database
}

func NewPaymentRepo(dbClient *db.Database) *PaymentRepoImpl {
	return &PaymentRepoImpl{
		dbClient: dbClient,
	}
}

func (repo PaymentRepoImpl) Create(payment *entities.Payment) error {
	tx := repo.dbClient.Core.Create(payment)
	return tx.Error
}

func (repo PaymentRepoImpl) FetchByAuthority(authority string) (*entities.Payment, error) {
	var payment *entities.Payment
	tx := repo.dbClient.Core.
		Where("authority = ?", authority).
		Preload("User").
		First(&payment)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return payment, nil
}

func (repo PaymentRepoImpl) Update(payment *entities.Payment) error {
	tx := repo.dbClient.Core.Updates(payment)
	return tx.Error
}

func (repo PaymentRepoImpl) FetchPageable(page, pageSize int) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	tx := repo.dbClient.Core.
		Offset(page * pageSize).
		Limit(pageSize).
		Order("created_at desc").
		Find(&payments)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return payments, nil
}

func (repo PaymentRepoImpl) FetchCount() (int, error) {
	var count int64
	tx := repo.dbClient.Core.
		Model(&entities.Payment{}).
		Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
