package repository

import (
	"errors"
	"fmt"
	"strconv"
	"tukangku/features/transaction"
	"tukangku/helper/midtrans"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	NoInvoice  string
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

	var id = strconv.Itoa(int(input.ID))
	input.NoInvoice = "TUKANGKU-ID-" + id

	midtrans := midtrans.MidtransCreateToken(int(input.ID), int(JobPrice))

	fmt.Println("Redirect URL:", midtrans.RedirectURL)
    fmt.Println("Token:", midtrans.Token)

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
	result.NoInvoice = input.NoInvoice

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

	// status := midtrans.MidtransStatus(int(transactionID))
	// transactions.Status = status

	// if err := ct.db.Save(&transactions).Error; err != nil {
	// 	return nil, err
	// }
	if transactions.ID == 0 {
		err := errors.New("no transactions")
		return nil, err
	}

	result := &transaction.Transaction{
		ID:         transactions.ID,
		JobID:      transactions.JobID,
		TotalPrice: transactions.TotalPrice,
		Status:     transactions.Status,
		Url:        transactions.Url,
		Token:      transactions.Token,
		NoInvoice:  transactions.NoInvoice,
	}

	return result, nil

}

func (cb *TransactionQuery) CallBack(noInvoice string) (*transaction.TransactionList, error) {
	var transactions Transaction
	if err := cb.db.Table("transactions").Where("no_invoice = ?", noInvoice).Find(&transactions).Error; err != nil {
		fmt.Println("transactions = ", transactions)
		return &transaction.TransactionList{}, err
	}

	if transactions.ID == 0 {
		err := errors.New("no transactions")
		return nil, err
	}

	ms := midtrans.MidtransStatus(noInvoice)
	transactions.Status = ms

	if err := cb.db.Save(&transactions).Error; err != nil {
		return nil, err
	}

	result := &transaction.TransactionList{
		ID:         transactions.ID,
		NoInvoice:  transactions.NoInvoice,
		JobID:      transactions.JobID,
		TotalPrice: transactions.TotalPrice,
		Status:     ms,
		Timestamp:  transactions.CreatedAt,
		Token:      transactions.Token,
		Url:        transactions.Url,
	}
	return result, nil

}
