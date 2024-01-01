package services

import (
	"errors"
	"tukangku/features/transaction"
	"tukangku/helper/jwt"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type TransactionService struct {
	repo transaction.Repository
}


func New(r transaction.Repository) transaction.Service {
	return &TransactionService{
		repo: r,
	}
}

func (at *TransactionService) AddTransaction(token *golangjwt.Token, JobID uint, JobPrice uint) (transaction.Transaction, error) {
	userID, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.Transaction{}, err
	}

	result, err := at.repo.AddTransaction(userID, JobID, JobPrice)

	return result, err
}

func (ct *TransactionService) CheckTransaction(transactionID uint) (transaction.Transaction, error) {
	result, err := ct.repo.CheckTransaction(transactionID)
	if err != nil {
		return transaction.Transaction{}, err
	}
	return *result, err
}


func (cb *TransactionService) CallBack(noInvoice string) (transaction.TransactionList, error) {
	result, err := cb.repo.CallBack(noInvoice)
	if err != nil {
		return transaction.TransactionList{}, err
	}
	if result == nil {
		return transaction.TransactionList{}, errors.New("result is nil")
	}
	return *result, err
}