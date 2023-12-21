package jobs

type CreateRequest struct {
	WorkerID   uint   `json:"worker_id"`
	ClientName string `json:"client_name"`
	Role       string `json:"role"`
	Category   string `json:"category"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Deskripsi  string `json:"deskripsi"`
}

type CreateResponse struct {
	ID   uint   `json:"job_id"`
	Foto string `json:"foto"`

	WorkerName string `json:"worker_name"`
	ClientName string `json:"client_name"`
	Category   string `json:"category"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Price      int    `json:"price"`
	Deskripsi  string `json:"deskripsi"`
	Status     string `json:"status"`
	Address    string `json:"address"`
}
type GetRequest struct {
	Role string `json:"role"`
}

type UpdateRequest struct {
	Price     int    `json:"price"`
	Deskripsi string `json:"deskripsi"`
	Status    string `json:"status"`
	Role      string `json:"role"`
}
