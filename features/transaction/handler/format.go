package handler

type TransactionReq struct {
	JobID    uint `json:"job_id"`
	JobPrice uint `json:"job_price"`
}

type TransactionRes struct {
	ID       uint `json:"transaction_id"`
	JobID    uint `json:"job_id"`
	JobPrice uint `json:"job_price"`
	Status   string `json:"status"`
	Token    string `json:"token"`
	Url      string `json:"url"`
}
