package model

import (
	"errors"
	"fmt"
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
	Foto      string
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

type NotifModel struct {
	gorm.Model
	UserID  uint `gorm:"not null"`
	Message string
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
	var worker = new(UserModel)
	result = jq.db.Where("id = ?", newJobs.WorkerID).First(&worker)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan worker")
	}
	input.Address = client.Alamat
	input.WorkerID = newJobs.WorkerID
	input.ClientID = newJobs.ClientID
	input.Category = newJobs.Category
	input.StartDate = newJobs.StartDate
	input.EndDate = newJobs.EndDate
	input.Foto = worker.Foto
	input.Price = 0
	input.Deskripsi = newJobs.Deskripsi
	input.Status = "pending"

	if err := jq.db.Create(&input).Error; err != nil {
		return jobs.Jobs{}, err
	}
	// bikin notif dulu
	var notif = new(NotifModel)
	notif.UserID = newJobs.WorkerID
	notif.Message = "Anda mendapatkan request job baru!"
	if err := jq.db.Create(&notif).Error; err != nil {
		return jobs.Jobs{}, err
	}

	// ngambil data dari repo untuk dikembalikan

	// fmt.Println(input.ID)
	// fmt.Println(worker)
	var response = jobs.Jobs{
		ID:         input.ID,
		Foto:       input.Foto,
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

func (jq *jobQuery) GetJobs(userID uint, role string) ([]jobs.Jobs, error) {
	var proses = new([]JobModel)
	switch role {
	case "worker":
		if err := jq.db.Where("worker_id = ?", userID).Order("created_at desc").Find(&proses).Error; err != nil {
			return nil, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, nil
		}

		var worker = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&worker)
		if result.Error != nil {
			return []jobs.Jobs{}, errors.New("tidak ditemukan worker, 404")
		}

		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var output = new(jobs.Jobs)
			var client = new(UserModel)
			result = jq.db.Where("id = ?", element.ClientID).First(&client)
			if result.Error != nil {
				return []jobs.Jobs{}, errors.New("tidak ditemukan client, 404")
			}
			output.ID = element.ID
			output.WorkerID = element.WorkerID
			output.WorkerName = worker.Nama
			output.ClientID = element.ClientID
			output.ClientName = client.Nama
			output.Category = element.Category
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price
			output.Deskripsi = element.Deskripsi
			output.Status = element.Status
			output.Address = element.Address
			output.Foto = element.Foto
			*outputs = append(*outputs, *output)
		}
		return *outputs, nil
	case "client":
		if err := jq.db.Where("client_id = ?", userID).Order("created_at desc").Find(&proses).Error; err != nil {
			return nil, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, nil
		}

		var client = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&client)
		if result.Error != nil {
			return []jobs.Jobs{}, errors.New("tidak ditemukan client, 404")
		}
		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var worker = new(UserModel)
			result = jq.db.Where("id = ?", element.WorkerID).First(&worker)
			if result.Error != nil {
				return []jobs.Jobs{}, errors.New("tidak ditemukan worker, 404")
			}
			var output = new(jobs.Jobs)
			output.ID = element.ID
			output.WorkerID = element.WorkerID
			output.WorkerName = worker.Nama
			output.ClientID = element.ClientID
			output.ClientName = client.Nama
			output.Category = element.Category
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price
			output.Deskripsi = element.Deskripsi
			output.Status = element.Status
			output.Address = element.Address
			output.Foto = element.Foto
			*outputs = append(*outputs, *output)
		}
		return *outputs, nil
	default:
		return nil, errors.New("kesalahan pada role")
	}

}
func (jq *jobQuery) GetJobsByStatus(userID uint, status string, role string) ([]jobs.Jobs, error) {
	var proses = new([]JobModel)
	switch role {
	case "worker":
		if err := jq.db.Where("worker_id = ? AND status = ?", userID, status).Order("created_at desc").Find(&proses).Error; err != nil {
			return nil, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, nil
		}
		var worker = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&worker)
		if result.Error != nil {
			return []jobs.Jobs{}, errors.New("tidak ditemukan worker, 404")
		}

		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var output = new(jobs.Jobs)
			var client = new(UserModel)
			result = jq.db.Where("id = ?", element.ClientID).First(&client)
			if result.Error != nil {
				return []jobs.Jobs{}, errors.New("tidak ditemukan client, 404")
			}
			output.ID = element.ID
			output.Foto = element.Foto
			output.WorkerID = element.WorkerID
			output.WorkerName = worker.Nama
			output.ClientID = element.ClientID
			output.ClientName = client.Nama
			output.Category = element.Category
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price
			output.Deskripsi = element.Deskripsi
			output.Status = element.Status
			output.Address = element.Address
			*outputs = append(*outputs, *output)
		}
		return *outputs, nil
	case "client":
		if err := jq.db.Where("client_id = ? AND status = ?", userID, status).Order("created_at desc").Find(&proses).Error; err != nil {
			return nil, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, nil
		}

		var client = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&client)
		if result.Error != nil {
			return []jobs.Jobs{}, errors.New("tidak ditemukan client, 404")
		}
		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var worker = new(UserModel)
			result = jq.db.Where("id = ?", element.WorkerID).First(&worker)
			if result.Error != nil {
				return []jobs.Jobs{}, errors.New("tidak ditemukan worker, 404")
			}
			var output = new(jobs.Jobs)
			output.ID = element.ID
			output.Foto = element.Foto
			output.WorkerID = element.WorkerID
			output.WorkerName = worker.Nama
			output.ClientID = element.ClientID
			output.ClientName = client.Nama
			output.Category = element.Category
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price
			output.Deskripsi = element.Deskripsi
			output.Status = element.Status
			output.Address = element.Address
			*outputs = append(*outputs, *output)
		}
		return *outputs, nil
	default:
		return nil, errors.New("kesalahan pada role")
	}
}

func (jq *jobQuery) GetJob(jobID uint) (jobs.Jobs, error) {
	var proses = new(JobModel)

	result := jq.db.Where("id = ?", jobID).First(&proses)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan jobs")
	}
	var output = new(jobs.Jobs)
	var client = new(UserModel)
	result = jq.db.Where("id = ?", proses.ClientID).First(&client)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan client, 404")
	}
	var worker = new(UserModel)
	result = jq.db.Where("id = ?", proses.ClientID).First(&worker)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan woker, 404")
	}
	output.ID = proses.ID
	output.Foto = proses.Foto
	output.WorkerID = proses.WorkerID
	output.WorkerName = worker.Nama
	output.ClientID = proses.ClientID
	output.ClientName = client.Nama
	output.Category = proses.Category
	output.StartDate = proses.StartDate
	output.EndDate = proses.EndDate
	output.Price = proses.Price
	output.Deskripsi = proses.Deskripsi
	output.Status = proses.Status
	output.Address = proses.Address
	fmt.Println(output, "repo")
	return *output, nil

}
func (jq *jobQuery) UpdateJob(update jobs.Jobs) (jobs.Jobs, error) {
	var proses = new(JobModel)
	result := jq.db.Where("id = ?", update.ID).First(&proses)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan jobs")
	}
	if proses.Status == "finished" {
		return jobs.Jobs{}, errors.New("jobs tidak boleh diubah, 403")
	}
	// cek id updater
	if update.Role == "client" {
		if update.ClientID != proses.ClientID {
			return jobs.Jobs{}, errors.New("id client tidak cocok, 403")
		}
	} else {
		if update.WorkerID != proses.WorkerID {
			return jobs.Jobs{}, errors.New("id worker tidak cocok, 403")
		}
	}
	if update.Price != 0 {
		proses.Price = update.Price

	}
	if update.Deskripsi != "" {
		proses.Deskripsi = update.Deskripsi
	}
	if update.Status != "" {
		proses.Status = update.Status
	}
	// proses
	result = jq.db.Save(&proses)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("eror saat menyimpan data, 500")
	}
	// bikin notif
	var notif = new(NotifModel)
	switch update.Role {
	case "client":
		var worker = new(UserModel)
		result = jq.db.Where("id = ?", proses.WorkerID).First(&worker)
		if result.Error != nil {
			return jobs.Jobs{}, errors.New("tidak ditemukan worker, notif, 404")
		}
		notif.UserID = update.ClientID
		notif.Message = fmt.Sprintf("Worker %v telah mengubah detail pada Job Request Anda, Job ID: %v", worker.Nama, update.ID)

		result = jq.db.Create(&notif)
		if result.Error != nil {
			return jobs.Jobs{}, errors.New("kesalahan saat membuat notif")
		}
	case "worker":
		var client = new(UserModel)
		result = jq.db.Where("id = ?", proses.ClientID).First(&client)
		if result.Error != nil {
			return jobs.Jobs{}, errors.New("tidak ditemukan client, notif, 404")
		}
		notif.UserID = update.WorkerID
		notif.Message = fmt.Sprintf("Request %v diubah oleh client: %v. Harap cek penawaran terkait.", update.ID, client.Nama)

		result = jq.db.Create(&notif)
		if result.Error != nil {
			return jobs.Jobs{}, errors.New("kesalahan saat membuat notif")
		}
	}
	var output = new(jobs.Jobs)
	var client = new(UserModel)
	result = jq.db.Where("id = ?", proses.ClientID).First(&client)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan client, 404")
	}
	var worker = new(UserModel)
	result = jq.db.Where("id = ?", proses.WorkerID).First(&worker)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan woker, 404")
	}
	output.ID = proses.ID
	output.Foto = proses.Foto
	output.WorkerID = proses.WorkerID
	output.WorkerName = worker.Nama
	output.ClientID = proses.ClientID
	output.ClientName = client.Nama
	output.Category = proses.Category
	output.StartDate = proses.StartDate
	output.EndDate = proses.EndDate
	output.Price = proses.Price
	output.Deskripsi = proses.Deskripsi
	output.Status = proses.Status
	output.Address = proses.Address

	return *output, nil

}
