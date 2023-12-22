package jobs

import "github.com/labstack/echo/v4"

type Jobs struct {
	ID         uint   `json:"job_id"`
	Foto       string `json:"foto"`
	WorkerName string `json:"worker_name"`
	WorkerID   uint   `json:"worker_id"`
	ClientName string `json:"client_name"`
	ClientID   uint   `json:"client_id"`
	Role       string `json:"role"`
	Category   string `json:"category"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Price      int    `json:"harga"`
	Deskripsi  string `json:"deskripsi"`
	Note       string `json:"note_negosiasi"`
	Status     string `json:"status"`
	Address    string `json:"alamat"`
}

type Handler interface {
	Create() echo.HandlerFunc
	GetJobs() echo.HandlerFunc
	GetJob() echo.HandlerFunc
	UpdateJob() echo.HandlerFunc
}

type Service interface {
	Create(newJobs Jobs) (Jobs, error)
	GetJobs(userID uint, status string, role string, page int, pagesize int) ([]Jobs, int, error)
	GetJob(jobID uint, role string) (Jobs, error)
	UpdateJob(update Jobs) (Jobs, error)
}

type Repository interface {
	Create(newJobs Jobs) (Jobs, error)
	GetJobs(userID uint, role string, page int, pagesize int) ([]Jobs, int, error)
	GetJobsByStatus(userID uint, status string, role string, page int, pagesize int) ([]Jobs, int, error)
	GetJob(jobID uint, role string) (Jobs, error)
	UpdateJob(update Jobs) (Jobs, error)
}
