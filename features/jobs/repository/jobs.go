package model

import (
	"errors"
	"tukangku/features/jobs"
	"tukangku/features/skill/repository"

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
	Address   string
}

type UserModel struct {
	gorm.Model
	Nama     string
	UserName string
	Password string
	Email    string
	NoHp     string
	Alamat   string
	Foto     string
	Role     string
	Skill    []repository.SkillModel `gorm:"many2many:user_skills;"`
	// Category []model.SkillModel `gorm:"foreignKey:Skill"`
	// SkillUser []skill.Skills `gorm:"foreignKey:Skill"`
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
	var input = new(JobModel)
	var client = new(UserModel)
	result := jq.db.Where("id = ?", newJobs.ClientID).First(&client)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan client")
	}
	input.Address = client.Alamat
	input.WorkerID = newJobs.WorkerID
	input.ClientID = newJobs.ClientID
	input.Category = newJobs.Category
	input.StartDate = newJobs.StartDate
	input.EndDate = newJobs.EndDate

	input.Price = 0
	input.Deskripsi = newJobs.Deskripsi
	input.Status = "pending"

	if err := jq.db.Create(&input).Error; err != nil {
		return jobs.Jobs{}, err
	}
	// bikin notif dulu

	// ngambil data dari repo untuk dikembalikan
	var worker = new(UserModel)
	result = jq.db.Where("id = ?", newJobs.WorkerID).First(&worker)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan worker")
	}
	// fmt.Println(input.ID)
	// fmt.Println(worker)
	var response = jobs.Jobs{
		ID:         input.ID,
		WorkerID:   input.WorkerID,
		WorkerName: worker.Nama,
		ClientID:   input.ClientID,
		Category:   input.Category,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		Price:      input.Price,
		Deskripsi:  input.Deskripsi,
		Status:     input.Status,
		Address:    input.Address,
	}
	// fmt.Println(response.ID)
	// fmt.Println(response.WorkerName)
	return response, nil
}
