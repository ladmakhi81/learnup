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
	unitOfWork db.UnitOfWork
}

func NewTransactionService(
	unitOfWork db.UnitOfWork,
) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		unitOfWork: unitOfWork,
	}
}

func (svc TransactionServiceImpl) FetchPageable(page, pageSize int) ([]*entities.Transaction, int, error) {
	const operationName = "TransactionServiceImpl.FetchPageable"
	transactions, count, err := svc.unitOfWork.TransactionRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset: &page,
			Limit:  &pageSize,
		},
	)
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching transactions",
			operationName,
			err,
		)
	}
	return transactions, count, nil
}
