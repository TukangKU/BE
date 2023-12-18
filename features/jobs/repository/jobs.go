package model

import (
	"tukangku/features/jobs"

	"gorm.io/gorm"
)

type JobModel struct {
	gorm.Model
	WorkerID  uint
	ClientID  uint
	Category  string
	StartDate string
	EndDate   string
	Price     int
	Deskripsi string
	Status    string
}

type jobQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) jobs.Repository {
	return &jobQuery{
		db: db,
	}
}

func (jq *jobQuery) Create(newJobs jobs.Jobs) (jobs.Jobs, error) {
	var input = JobModel{
		WorkerID:  newJobs.WorkerID,
		ClientID:  newJobs.ClientID,
		Category:  newJobs.Category,
		StartDate: newJobs.StartDate,
		EndDate:   newJobs.EndDate,

		Price:     0,
		Deskripsi: newJobs.Deskripsi,
		Status:    "pending",
	}

	if err := jq.db.Create(&input).Error; err != nil {
		return jobs.Jobs{}, err
	}
	// bikin notif dulu

	// ngambil data dari repo untuk dikembalikan
	newJobs.ID = input.ID

	return newJobs, nil
}
