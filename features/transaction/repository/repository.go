package repository

import (
	"errors"
	"tukangku/features/transaction"
	"tukangku/helper/midtrans"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	JobID      uint
	UserID     uint
	TotalPrice uint
	Status     string
	Token      string
	Url        string
}

type TransactionQuery struct {
	db *gorm.DB
}


func New(db *gorm.DB) transaction.Repository {
	return &TransactionQuery{
		db: db,
	}
}

func (at *TransactionQuery) AddTransaction(userID uint, JobID uint, JobPrice uint) (transaction.Transaction, error) {
	var input Transaction
	input.UserID = userID
	input.JobID = JobID
	input.TotalPrice = JobPrice
	input.Status = "pending"
	if err := at.db.Create(&input).Error; err != nil {
		return transaction.Transaction{}, err
	}

	midtrans := midtrans.MidtransCreateToken(int(input.ID), int(JobPrice))

	input.Url = midtrans.RedirectURL
	input.Token = midtrans.Token
	if err := at.db.Save(&input).Error; err != nil {
		return transaction.Transaction{}, err
	}

	var result transaction.Transaction
	result.ID = input.ID
	result.JobID = input.JobID
	result.TotalPrice = input.TotalPrice
	result.Status = input.Status
	result.Url = midtrans.RedirectURL
	result.Token = midtrans.Token

	return result, nil

}



// CheckTransaction implements transaction.Repository.
func (ct *TransactionQuery) CheckTransaction(transactionID uint) (*transaction.Transaction, error) {
	var transactions Transaction
	if err := ct.db.First(&transactions, transactionID).Error; err != nil {
		return nil, err
	}

	if transactions.ID == 0 {
		err := errors.New("transaction doesnt exist")
		return nil, err
	}

	status := midtrans.MidtransStatus(int(transactionID))
	transactions.Status = status

	if err := ct.db.Save(&transactions).Error; err != nil {
		return nil, err
	}

	result := &transaction.Transaction{
		ID: transactions.ID,
		JobID: transactions.JobID,
		TotalPrice: transactions.TotalPrice,
		Status: status,
	}

	return result, nil

}