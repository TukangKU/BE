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
		response.ID = result.ID
		response.WorkerName = result.WorkerName
		response.ClientName = input.ClientName
		response.Price = result.Price
		response.Category = result.Category
		response.StartDate = result.StartDate
		response.EndDate = result.EndDate
		response.Deskripsi = result.Deskripsi
		response.Status = result.Status
		response.Address = result.Address
		// fmt.Println(result, "handler")
		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)

	}
}

// get jobs with and without query
func (jc *jobsController) GetJobs() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(GetRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang di berikan tidak sesuai",
			})
		}
		userID, err := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusUnauthorized
			var message = "harap login"

			return responses.PrintResponse(c, statusCode, message, nil)
		}
		status := c.QueryParams().Get("status")
		result, err := jc.srv.GetJobs(userID, status, input.Role)
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusUnauthorized
			var message = "harap login"

			return responses.PrintResponse(c, statusCode, message, nil)
		}
		var statusCode = http.StatusOK
		var message = "sukses"
		var respon = new([]CreateResponse)
		for _, element := range result {
			var response = new(CreateResponse)
			response.ID = element.ID
			response.WorkerName = element.WorkerName
			response.ClientName = element.ClientName
			response.Price = element.Price
			response.Category = element.Category
			response.StartDate = element.StartDate
			response.EndDate = element.EndDate
			response.Deskripsi = element.Deskripsi
			response.Status = element.Status
			response.Address = element.Address
			*respon = append(*respon, *response)
		}

		return responses.PrintResponse(c, statusCode, message, respon)
	}
}
