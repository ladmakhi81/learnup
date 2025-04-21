package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	transactionDtoReq "github.com/ladmakhi81/learnup/internals/transaction/dto/req"
	transactionRepository "github.com/ladmakhi81/learnup/internals/transaction/repo"
	"github.com/ladmakhi81/learnup/types"
)

type TransactionService interface {
	Create(dto transactionDtoReq.CreateTransactionReq) (*entities.Transaction, error)
	FetchPageable(page, pageSize int) ([]*entities.Transaction, error)
	FetchCount() (int, error)
}

type TransactionServiceImpl struct {
	repo transactionRepository.TransactionRepo
}

func NewTransactionService(
	repo transactionRepository.TransactionRepo,
) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		repo: repo,
	}
}

func (svc TransactionServiceImpl) Create(dto transactionDtoReq.CreateTransactionReq) (*entities.Transaction, error) {
	transaction := &entities.Transaction{
		Amount:   dto.Amount,
		Type:     dto.Type,
		User:     dto.User,
		Phone:    dto.Phone,
		Tag:      dto.Tag,
		Currency: dto.Currency,
	}
	if err := svc.repo.Create(transaction); err != nil {
		return nil, types.NewServerError(
			"Error in creating transaction",
			"TransactionServiceImpl.Create",
			err,
		)
	}
	return transaction, nil
}

func (svc TransactionServiceImpl) FetchPageable(page, pageSize int) ([]*entities.Transaction, error) {
	transactions, transactionsErr := svc.repo.FetchPageable(page, pageSize)
	if transactionsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching transactions",
			"TransactionServiceImpl.FetchPageable",
			transactionsErr,
		)
	}
	return transactions, nil
}

func (svc TransactionServiceImpl) FetchCount() (int, error) {
	count, countErr := svc.repo.FetchCount()
	if countErr != nil {
		return 0, types.NewServerError(
			"Error in fetching count of transactions",
			"TransactionServiceImpl.FetchCount",
			countErr,
		)
	}
	return count, nil
}
