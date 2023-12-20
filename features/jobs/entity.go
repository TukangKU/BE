package jobs

import "github.com/labstack/echo/v4"

type Jobs struct {
	ID         uint   `json:"job_id"`
	WorkerName string `json:"worker_name"`
	WorkerID   uint   `json:"worker_id"`
	ClientName string `json:"client_name"`
	ClientID   uint   `json:"client_id"`
	Role       string `json:"role"`
	Category   string `json:"category"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Price      int    `json:"price"`
	Deskripsi  string `json:"deskripsi"`
	Status     string `json:"status"`
	Address    string `json:"address"`
}

type Handler interface {
	Create() echo.HandlerFunc
}

type Service interface {
	Create(newJobs Jobs) (Jobs, error)
}

type Repository interface {
	Create(newJobs Jobs) (Jobs, error)
}
