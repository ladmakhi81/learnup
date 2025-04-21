package repo

import (
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
)

type TransactionRepo interface {
	Create(transaction *entities.Transaction) error
	FetchPageable(page, pageSize int) ([]*entities.Transaction, error)
	FetchCount() (int, error)
}

type TransactionRepoImpl struct {
	dbClient *db.Database
}

func NewTransactionRepo(dbClient *db.Database) *TransactionRepoImpl {
	return &TransactionRepoImpl{
		dbClient: dbClient,
	}
}

func (repo TransactionRepoImpl) Create(transaction *entities.Transaction) error {
	tx := repo.dbClient.Core.Create(transaction)
	return tx.Error
}

func (repo TransactionRepoImpl) FetchPageable(page, pageSize int) ([]*entities.Transaction, error) {
	var transactions []*entities.Transaction
	tx := repo.dbClient.Core.
		Offset(page * pageSize).
		Limit(pageSize).
		Order("created_at desc").
		Find(&transactions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return transactions, nil
}

func (repo TransactionRepoImpl) FetchCount() (int, error) {
	var count int64
	tx := repo.dbClient.Core.
		Model(&entities.Transaction{}).
		Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
