package jobs

type CreateRequest struct {
	WorkerID  uint   `json:"worker_id"`
	SkillID   uint   `json:"skill_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Deskripsi string `json:"deskripsi"`
	Address   string `json:"alamat"`
}

type CreateResponse struct {
	ID   uint   `json:"job_id"`
	Foto string `json:"foto"`

	WorkerName string `json:"worker_name"`

	Category  string `json:"category"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Price     int    `json:"harga"`
	Deskripsi string `json:"deskripsi"`
	Status    string `json:"status"`
	Address   string `json:"alamat"`
}

type UpdateRequest struct {
	Price    int    `json:"harga"`
	NoteNego string `json:"note_negosiasi"`
	Status   string `json:"status"`
}

type GetJobsResponse struct {
	ID         uint   `json:"job_id"`
	ClientName string `json:"client_name"`
	WorkerName string `json:"worker_name"`
	Foto       string `json:"foto"`
	Category   string `json:"category"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Price      int    `json:"harga"`
	Status     string `json:"status"`
}

type GetJobResponse struct {
	ID         uint   `json:"job_id"`
	Category   string `json:"category"`
	WorkerName string `json:"worker_name"`
	ClientName string `json:"client_name"`
	Foto       string `json:"foto"`

	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Address   string `json:"alamat"`
	Price     int    `json:"harga"`
	Deskripsi string `json:"deskripsi"`
	Note      string `json:"note_negosiasi"`
	Status    string `json:"status"`
}
