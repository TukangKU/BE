package jobs

import (
	"net/http"
	"tukangku/features/jobs"
	"tukangku/helper/jwt"
	"tukangku/helper/responses"

	golangjwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type jobsController struct {
	srv jobs.Service
}

func New(s jobs.Service) jobs.Handler {
	return &jobsController{
		srv: s,
	}
}

// create jobs
func (jc *jobsController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		var input = new(CreateRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang di berikan tidak sesuai",
			})
		}

		var inputProcess = new(jobs.Jobs)
		inputProcess.ClientID = userID
		inputProcess.WorkerID = input.WorkerID
		inputProcess.Role = input.Role
		inputProcess.Category = input.Category
		inputProcess.StartDate = input.StartDate
		inputProcess.EndDate = input.EndDate
		inputProcess.Deskripsi = input.Deskripsi

		result, err := jc.srv.Create(*inputProcess)
		// error nya belum
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		var response = new(CreateResponse)
		response.WorkerName = result.WorkerName
		response.ClientName = input.ClientName
		response.Price = result.Price
		response.Category = result.Category
		response.StartDate = result.StartDate
		response.EndDate = result.EndDate
		response.Deskripsi = result.Deskripsi
		response.Status = result.Status

		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)

	}
}
