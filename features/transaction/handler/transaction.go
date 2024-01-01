package handler

import (
	"net/http"
	"strconv"
	"strings"
	"tukangku/features/transaction"
	"tukangku/helper/responses"

	gojwt "github.com/golang-jwt/jwt/v5"

	echo "github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	s transaction.Service
}

func New(s transaction.Service) transaction.Handler {
	return &TransactionHandler{
		s: s,
	}
}

func (at *TransactionHandler) AddTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(TransactionReq)
		if err := c.Bind(input); err != nil {
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"input yang diberikan tidak sesuai",
				nil)
		}
		result, err := at.s.AddTransaction(c.Get("user").(*gojwt.Token), input.JobID, input.JobPrice)

		if err != nil {
			c.Logger().Error("terjadi kesalahan", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return responses.PrintResponse(
					c, http.StatusBadRequest,
					"terjadi kesalahan input data",
					nil)
			}
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"transaction duplicate",
				nil)
		}

		var response = new(TransactionRes)
		response.ID = result.ID
		response.NoInvoice = result.NoInvoice
		response.JobID = result.JobID
		response.JobPrice = result.TotalPrice
		response.Status = result.Status
		response.Url = result.Url
		response.Token = result.Token

		return responses.PrintResponse(
			c, http.StatusCreated,
			"transaction created successfully",
			response)

	}
}

func (ct *TransactionHandler) CheckTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"id tidak valid",
				nil)
		}
		result, err := ct.s.CheckTransaction(uint(transactionID))

		if err != nil {
			c.Logger().Error("Error fetching : ", err.Error())
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"failed to retrieve data",
				nil)
		}

		var response = new(TransactionRes)
		response.ID = result.ID
		response.NoInvoice = result.NoInvoice
		response.JobID = result.JobID
		response.JobPrice = result.TotalPrice
		response.Status = result.Status

		return responses.PrintResponse(
			c, http.StatusOK,
			"transaction detail",
			response)
	}
}

func (cb *TransactionHandler) CallBack() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(CallBack)
		if err := c.Bind(input); err != nil {

			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"input tidak sesuai",
				nil)
		}
		result, err := cb.s.CallBack(input.OrderID)
		if err != nil {
			c.Logger().Error("something wrong: ", err.Error())
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"something wrong",
				nil)
		}
		return responses.PrintResponse(
			c, http.StatusOK,
			"Midtrans Callback",
			result)
	}
}
