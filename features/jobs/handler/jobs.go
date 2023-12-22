package jobs

import (
	"fmt"
	"math"
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

		workerID, err := strconv.Atoi(c.Param("worker_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "masukkan worker_id",
			})
		}
		userID, _ := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		userRole, _ := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))
		var input = new(CreateRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang di berikan tidak sesuai",
			})
		}

		var inputProcess = new(jobs.Jobs)
		inputProcess.ClientID = userID
		inputProcess.WorkerID = uint(workerID)
		inputProcess.Role = userRole
		inputProcess.Category = input.Category
		inputProcess.StartDate = input.StartDate
		inputProcess.EndDate = input.EndDate
		inputProcess.Deskripsi = input.Deskripsi
		inputProcess.Address = input.Address

		result, err := jc.srv.Create(*inputProcess)
		// error nya belum
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.Logger().Error("ERROR Register, explain:", err.Error())
				var statusCode = http.StatusNotFound
				var message = "data worker atau client tidak ditemukan"

				return responses.PrintResponse(c, statusCode, message, nil)
			}
			if strings.Contains(err.Error(), "bukan client") {
				c.Logger().Error("ERROR Register, explain:", err.Error())
				var statusCode = http.StatusUnauthorized
				var message = "anda bukan client"

				return responses.PrintResponse(c, statusCode, message, nil)
			}
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		var response = new(CreateResponse)
		response.ID = result.ID
		response.Foto = result.Foto
		response.WorkerName = result.WorkerName

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
		// get role
		userRole, err := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusUnauthorized
			var message = "harap login"

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		// get uid
		userID, err := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			if strings.Contains(err.Error(), "tidak ditemukan") {
				var statusCode = http.StatusNotFound
				var message = "tidak ditemukan"

				return responses.PrintResponse(c, statusCode, message, nil)
			}
			var statusCode = http.StatusUnauthorized
			var message = "harap login"

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		// get queries
		status := c.QueryParams().Get("status")
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.QueryParam("pagesize"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}

		// proses
		result, count, err := jc.srv.GetJobs(userID, status, userRole, page, pageSize)
		if err != nil {
			c.Logger().Error("ERROR Database, explain:", err.Error())
			if strings.Contains(err.Error(), "tidak ditemukan") {
				var statusCode = http.StatusNotFound
				var message = "ada  data yang hilang atau terhapus"

				return responses.PrintResponse(c, statusCode, message, nil)
			} else if strings.Contains(err.Error(), "token") {
				var statusCode = http.StatusUnauthorized
				var message = "salah token"

				return responses.PrintResponse(c, statusCode, message, nil)
			}
			var statusCode = http.StatusUnauthorized
			var message = "harap login"

			return responses.PrintResponse(c, statusCode, message, nil)
		}
		totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
		if page > totalPages {
			var statusCode = http.StatusNotFound
			var message = "index out of bounds"

			return responses.PrintResponse(c, statusCode, message, nil)
		}
		// proses response

		var respon = new([]GetJobsResponse)
		for _, element := range result {
			var response = new(GetJobsResponse)

			response.Foto = element.Foto
			response.WorkerName = element.WorkerName
			response.ClientName = element.ClientName

			response.Category = element.Category
			response.StartDate = element.StartDate
			response.EndDate = element.EndDate

			response.Status = element.Status

			*respon = append(*respon, *response)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get data",
			"data":    respon,
			"pagination": map[string]interface{}{
				"page":       page,
				"pagesize":   pageSize,
				"totalPages": totalPages},
		})
	}
}

func (jc *jobsController) GetJob() echo.HandlerFunc {
	return func(c echo.Context) error {

		userRole, err := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusUnauthorized
			var message = "harap login"

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		jobID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID tidak valid",
			})
		}
		result, err := jc.srv.GetJob(uint(jobID), userRole)

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
		var response = new(GetJobResponse)

		response.ID = result.ID
		response.Category = result.Category
		response.WorkerName = result.WorkerName
		response.ClientName = result.ClientName
		response.Foto = result.Foto
		response.StartDate = result.StartDate
		response.EndDate = result.EndDate
		response.Address = result.Address
		response.Price = result.Price

		response.Deskripsi = result.Deskripsi
		response.Note = result.Note
		response.Status = result.Status

		return responses.PrintResponse(c, http.StatusOK, "success create data", response)
	}
}

func (jc *jobsController) UpdateJob() echo.HandlerFunc {
	return func(c echo.Context) error {

		jobID, err := strconv.Atoi(c.Param("id"))
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
		userRole, _ := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "harap login",
			})
		}
		var proses = new(jobs.Jobs)
		switch userRole {
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
		proses.Note = request.NoteNego
		proses.Status = request.Status
		proses.ID = uint(jobID)
		proses.Role = userRole
		result, err := jc.srv.UpdateJob(*proses)

		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				fmt.Println(err.Error())
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"message": "data tidak ditemukan",
				})
			}
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "ada masalah di server",
			})
		}

		var response = new(GetJobResponse)
		response.ID = result.ID
		response.Category = result.Category
		response.WorkerName = result.WorkerName
		response.ClientName = result.ClientName
		response.Foto = result.Foto
		response.StartDate = result.StartDate
		response.EndDate = result.EndDate
		response.Address = result.Address
		response.Price = result.Price
		response.Deskripsi = result.Deskripsi
		response.Note = result.Note
		response.Status = result.Status
		// fmt.Println(result, "handler")
		return responses.PrintResponse(c, http.StatusOK, "success create data", response)

	}
}
