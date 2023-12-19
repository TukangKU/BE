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
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}
		result, err := at.s.AddTransaction(c.Get("user").(*gojwt.Token), input.JobID, input.JobPrice)

		if err != nil {
			c.Logger().Error("terjadi kesalahan", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "dobel input nama",
				})
			}
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "transaction duplicate",
			})
		}

		var response = new(TransactionRes)
		response.ID = result.ID
		response.JobID = result.JobID
		response.JobPrice = result.TotalPrice
		response.Status = result.Status
		response.Url = result.Url
		response.Token = result.Token
		response.NoInvoice = result.NoInvoice

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Transaction created successfully",
			"data":    response,
		})

	}
}

func (ct *TransactionHandler) CheckTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID tidak valid",
			})
		}
		result, err := ct.s.CheckTransaction(uint(transactionID))

		if err != nil {
			c.Logger().Error("Error fetching : ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve data",
			})
		}

		var response = new(TransactionRes)
		response.ID = result.ID
		response.NoInvoice = result.NoInvoice
		response.JobID = result.JobID
		response.JobPrice = result.TotalPrice
		response.Status = result.Status

		return c.JSON(http.StatusOK, map[string]any{
			"message": "Transaction Detail",
			"data":    response,
		})
	}
}



func (cb *TransactionHandler) CallBack() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(CallBack)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input tidak sesuai",
			})
		}
		result, err := cb.s.CallBack(input.NoInvoice)
		if err != nil {
			c.Logger().Error("something wrong: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "something wrong",
			})
		}
		return responses.PrintResponse(c, http.StatusOK, "Midtrans Callback", result)
	}
}
