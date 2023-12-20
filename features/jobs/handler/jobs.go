package jobs

import (
	"net/http"
	"strconv"
	"strings"
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

func (jc *jobsController) GetJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID tidak valid",
			})
		}
		result, err := jc.srv.GetJob(uint(jobID))
		c.Logger().Error("ERROR GetByID, explain:", err.Error())
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"message": "Job not found",
				})
			}

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error retrieving Job by ID",
			})
		}

		// respons
		var response = new(CreateResponse)

		response.ID = result.ID
		response.WorkerName = result.WorkerName
		response.ClientName = result.ClientName
		response.Price = result.Price
		response.Category = result.Category
		response.StartDate = result.StartDate
		response.EndDate = result.EndDate
		response.Deskripsi = result.Deskripsi
		response.Status = result.Status
		response.Address = result.Address
		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)

	}
}

func (jc *jobsController) UpdateJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID tidak valid",
			})
		}

		var request = new(UpdateRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang di berikan tidak sesuai",
			})
		}
		userID, err := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "harap login",
			})
		}
		var proses = new(jobs.Jobs)
		switch proses.Role {
		case "client":
			proses.ClientID = userID
		case "worker":
			proses.WorkerID = userID
		default:
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "role tidak dikenali",
			})
		}
		proses.Price = request.Price
		proses.Deskripsi = request.Deskripsi
		proses.Status = request.Status
		proses.ID = uint(jobID)
		proses.Role = request.Role
		result, err := jc.srv.UpdateJob(*proses)

		if err != nil {

			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "eror belum disetting handler ke servis",
			})
		}

		var response = new(CreateResponse)
		response.ID = result.ID
		response.WorkerName = result.WorkerName
		response.ClientName = result.ClientName
		response.Price = result.Price
		response.Category = result.Category
		response.StartDate = result.StartDate
		response.EndDate = result.EndDate
		response.Deskripsi = result.Deskripsi
		response.Status = result.Status
		response.Address = result.Address
		// fmt.Println(result, "handler")
		return responses.PrintResponse(c, http.StatusOK, "success create data", response)

	}
}
