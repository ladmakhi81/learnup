package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/types"
)

type TransactionService interface {
	FetchPageable(page, pageSize int) ([]*entities.Transaction, int, error)
}

type TransactionServiceImpl struct {
	repo *db.Repositories
}

func NewTransactionService(
	repo *db.Repositories,
) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		repo: repo,
	}
}

func (svc TransactionServiceImpl) FetchPageable(page, pageSize int) ([]*entities.Transaction, int, error) {
	transactions, count, transactionsErr := svc.repo.TransactionRepo.GetPaginated(
		repositories.GetPaginatedOptions{
			Offset: &page,
			Limit:  &pageSize,
		},
	)
	if transactionsErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching transactions",
			"TransactionServiceImpl.FetchPageable",
			transactionsErr,
		)
	}
	return transactions, count, nil
}
