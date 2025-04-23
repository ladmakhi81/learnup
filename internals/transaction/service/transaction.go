package service

import (
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/db/repositories"
	"github.com/ladmakhi81/learnup/shared/types"
)

type TransactionService interface {
	FetchPageable(page, pageSize int) ([]*entities.Transaction, int, error)
}

type transactionService struct {
	unitOfWork db.UnitOfWork
}

func NewTransactionSvc(unitOfWork db.UnitOfWork) TransactionService {
	return &transactionService{unitOfWork: unitOfWork}
}

func (svc transactionService) FetchPageable(page, pageSize int) ([]*entities.Transaction, int, error) {
	const operationName = "transactionService.FetchPageable"
	transactions, count, err := svc.unitOfWork.TransactionRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset: &page,
			Limit:  &pageSize,
		},
	)
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching transactions", operationName, err)
	}
	return transactions, count, nil
}
