package transaction

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	ID         uint
	JobID      uint
	TotalPrice uint
	Status     string
	Token      string
	Url        string
}

type TransactionList struct {
	ID         uint
	JobID      uint
	TotalPrice uint
	Status     string
	Token      string
	Url        string
	Timestamp  time.Time `json:"timestamp"`
}

type Handler interface {
	AddTransaction() echo.HandlerFunc
	CheckTransaction() echo.HandlerFunc
}

type Repository interface {
	AddTransaction(userID uint, JobID uint, JobPrice uint) (Transaction, error)
	CheckTransaction(transactionID uint) (*Transaction, error)
}

type Service interface {
	AddTransaction(token *jwt.Token, JobID uint, JobPrice uint) (Transaction, error)
	CheckTransaction(transactionID uint) (Transaction, error)
}
