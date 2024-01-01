package services_test

import (
	"errors"
	"testing"
	"tukangku/features/transaction"
	"tukangku/features/transaction/mocks"
	"tukangku/features/transaction/services"

	"github.com/stretchr/testify/assert"

	"tukangku/helper/jwt"

	gojwt "github.com/golang-jwt/jwt/v5"
)

var userID = uint(1)
var role = "worker"
var str, _ = jwt.GenerateJWT(userID, role)
var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyy!!!"), nil
})
var invalidToken, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3xxx!!!"), nil
})

func TestAddTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	t.Run("Success", func(t *testing.T) {
		mockTransaction := transaction.Transaction{
			ID:         1,
			JobID:      2,
			NoInvoice:  "TUKANGKU-01",
			TotalPrice: 10000,
		}

		repo.On("AddTransaction", userID, mockTransaction.JobID, mockTransaction.TotalPrice).Return(mockTransaction, nil).Once()

		result, err := s.AddTransaction(token, mockTransaction.JobID, mockTransaction.TotalPrice)

		repo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, mockTransaction, result)
	})

	t.Run("InvalidToken", func(t *testing.T) {

		mockTransaction := transaction.Transaction{
			ID:         1,
			JobID:      2,
			NoInvoice:  "TUKANGKU-01",
			TotalPrice: 10000,
		}

		result, err := s.AddTransaction(invalidToken, mockTransaction.JobID, mockTransaction.TotalPrice)

		assert.Error(t, err)
		assert.Equal(t, transaction.Transaction{}, result)

	})

}

func TestCheckTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	service := services.New(repo)

	transactionID := uint(1)

	t.Run("Success", func(t *testing.T) {
		expectedTransaction := transaction.Transaction{ID: 1, JobID: 1, TotalPrice: 10000, Url: ""}
		repo.On("CheckTransaction", transactionID).Return(&expectedTransaction, nil).Once()

		result, err := service.CheckTransaction(transactionID)

		repo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, result)

	})

	t.Run("RepositoryError", func(t *testing.T) {
		repo.On("CheckTransaction", uint(0)).Return(nil, errors.New("repository error")).Once()

		_, err := service.CheckTransaction(uint(0))

		repo.AssertExpectations(t)

		assert.Error(t, err)
	})

	t.Run("TransactionNotFound", func(t *testing.T) {
		repo.On("CheckTransaction", transactionID).Return(nil, errors.New("transaction not found")).Once()

		_, err := service.CheckTransaction(transactionID)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, "transaction not found", err.Error())
	})
}

func TestCallBack(t *testing.T) {
	t.Run("ErrorFromRepository", func(t *testing.T) {
		repo := new(mocks.Repository)

		service := services.New(repo)

		repo.On("CallBack", "invoice123").Return(nil, errors.New("callback error"))

		result, err := service.CallBack("invoice123")

		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionList{}, result)
		assert.Equal(t, "callback error", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("ResultIsNil", func(t *testing.T) {
		repo := new(mocks.Repository)

		service := services.New(repo)

		repo.On("CallBack", "invoice123").Return(nil, nil)

		result, err := service.CallBack("invoice123")

		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionList{}, result)
		assert.Equal(t, "result is nil", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("SuccessCase", func(t *testing.T) {
		repo := new(mocks.Repository)

		service := services.New(repo)

		expectedResult := &transaction.TransactionList{}
		repo.On("CallBack", "invoice123").Return(expectedResult, nil)

		result, err := service.CallBack("invoice123")

		assert.NoError(t, err)
		assert.Equal(t, *expectedResult, result)

		repo.AssertExpectations(t)
	})
}
