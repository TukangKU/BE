package model

import (
	"errors"
	"strings"
	"time"
	"tukangku/features/jobs"

	"gorm.io/gorm"
)

type JobModel struct {
	gorm.Model
	WorkerID      uint       `gorm:"not null"`
	ClientID      uint       `gorm:"not null"`
	Client        UserModel  `gorm:"foreignKey:ClientID;"`
	Category      uint       `gorm:"not null"`
	CategoryModel SkillModel `gorm:"foreignKey:Category;"`
	StartDate     string     `gorm:"not null"`
	EndDate       string     `gorm:"not null"`
	Price         int
	Deskripsi     string
	Status        string
	Address       string
	NoteNego      string
	StatusPay     Transaction `gorm:"foreignKey:Status;"`
}

type UserModel struct {
	gorm.Model
	Nama     string
	UserName string `gorm:"unique"`
	Password string
	Email    string `json:"email" gorm:"unique"`
	NoHp     string
	Alamat   string
	Foto     string
	Role     string
	Skill    []SkillModel `gorm:"many2many:user_skills;"`
	Jobs     []JobModel   `gorm:"foreignKey:WorkerID"`
	Requests []JobModel   `gorm:"foreignKey:ClientID"`
}

type SkillModel struct {
	ID        uint   `gorm:"primarykey"`
	NamaSkill string `json:"skill"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Jobs      []JobModel     `gorm:"foreignKey:Category"`
}

type Transaction struct {
	gorm.Model
	NoInvoice  string
	JobID      uint
	UserID     uint
	TotalPrice uint
	Status     string
	Token      string
	Url        string
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
	if newJobs.Role != "client" {
		return jobs.Jobs{}, errors.New("anda bukan client")
	}
	result := jq.db.Where("id = ?", newJobs.ClientID).First(&client)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan client")
	}
	var worker = new(UserModel)
	result = jq.db.Where("id = ?", newJobs.WorkerID).First(&worker)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan worker")
	}
	input.Address = newJobs.Address
	input.WorkerID = newJobs.WorkerID
	input.ClientID = newJobs.ClientID
	input.Category = newJobs.SkillID
	input.StartDate = newJobs.StartDate
	input.EndDate = newJobs.EndDate

	input.Price = 0
	input.Deskripsi = newJobs.Deskripsi
	input.Status = "pending"

	if err := jq.db.Create(&input).Error; err != nil {
		return jobs.Jobs{}, err
	}

	var skill = new(SkillModel)
	result = jq.db.Where("id = ?", newJobs.SkillID).First(&skill)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan client")
	}

	var response = jobs.Jobs{
		ID:         input.ID,
		Foto:       worker.Foto,
		NoHp:       worker.NoHp,
		WorkerName: worker.Nama,
		ClientID:   input.ClientID,
		Category:   skill.NamaSkill,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		Price:      input.Price,
		Deskripsi:  input.Deskripsi,
		Status:     input.Status,
		Address:    input.Address,
	}

	return response, nil
}

func (jq *jobQuery) GetJobs(userID uint, role string, page int, pagesize int) ([]jobs.Jobs, int, error) {
	var proses = new([]JobModel)
	var prePagination = new([]JobModel)
	var totalCount int64
	offset := (page - 1) * pagesize

	switch role {
	case "worker":
		var worker = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&worker)
		if result.Error != nil {
			return []jobs.Jobs{}, 0, errors.New("tidak ditemukan worker, 404")
		}
		if role != worker.Role {
			return []jobs.Jobs{}, 0, errors.New("sepertinya anda salah memasukkan token")
		}
		if err := jq.db.
			Where("worker_id = ?", userID).Order("updated_at desc").
			Find(&prePagination).
			Count(&totalCount).Error; err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil, 0, nil
			}
			return nil, 0, errors.New("server error")
		}

		if err := jq.db.
			Where("worker_id = ?", userID).Order("updated_at desc").
			Offset(offset).
			Limit(pagesize).
			Find(&proses).Error; err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil, 0, nil
			}
			return nil, 0, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, int(totalCount), nil
		}

		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var output = new(jobs.Jobs)
			var client = new(UserModel)
			result = jq.db.Where("id = ?", element.ClientID).First(&client)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client, 404")
			}
			output.WorkerName = worker.Nama
			output.ID = element.ID
			output.ClientName = client.Nama
			output.Foto = client.Foto

			var skill = new(SkillModel)
			result = jq.db.Where("id = ?", element.Category).First(&skill)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client")
			}
			output.Category = skill.NamaSkill
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price

			output.Status = element.Status

			*outputs = append(*outputs, *output)
		}
		return *outputs, int(totalCount), nil
	case "client":

		var client = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&client)
		if result.Error != nil {
			return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client, 404")
		}
		if role != client.Role {
			return []jobs.Jobs{}, 0, errors.New("sepertinya anda salah memasukkan token")
		}

		if err := jq.db.
			Where("client_id = ?", userID).Order("updated_at desc").
			Find(&prePagination).
			Count(&totalCount).Error; err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil, 0, nil
			}
			return nil, 0, errors.New("server error")
		}
		if err := jq.db.Where("client_id = ?", userID).Order("updated_at desc").
			Offset(offset).
			Limit(pagesize).
			Find(&proses).Error; err != nil {
			return nil, 0, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, int(totalCount), nil
		}

		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var worker = new(UserModel)
			result = jq.db.Where("id = ?", element.WorkerID).First(&worker)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan worker, 404")
			}
			var output = new(jobs.Jobs)
			output.ID = element.ID
			output.WorkerName = worker.Nama

			output.ClientName = client.Nama
			output.Foto = worker.Foto

			var skill = new(SkillModel)
			result = jq.db.Where("id = ?", element.Category).First(&skill)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client")
			}
			output.Category = skill.NamaSkill
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price

			output.Status = element.Status

			*outputs = append(*outputs, *output)
		}
		return *outputs, int(totalCount), nil
	default:
		return nil, 0, errors.New("kesalahan pada role")
	}

}
func (jq *jobQuery) GetJobsByStatus(userID uint, status string, role string, page int, pagesize int) ([]jobs.Jobs, int, error) {
	var proses = new([]JobModel)
	var prePagination = new([]JobModel)
	var totalCount int64
	offset := (page - 1) * pagesize
	switch role {
	case "worker":
		var worker = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&worker)
		if result.Error != nil {
			return []jobs.Jobs{}, 0, errors.New("tidak ditemukan worker, 404")
		}
		if role != worker.Role {
			return []jobs.Jobs{}, 0, errors.New("sepertinya anda salah memasukkan token")
		}
		if err := jq.db.
			Where("worker_id = ? AND status = ?", userID, status).
			Find(&prePagination).
			Count(&totalCount).Error; err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil, 0, nil
			}
			return nil, 0, errors.New("server error")
		}
		if err := jq.db.
			Where("worker_id = ? AND status = ?", userID, status).
			Order("updated_at desc").
			Offset(offset).
			Limit(pagesize).
			Find(&proses).
			Error; err != nil {
			return nil, 0, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, int(totalCount), nil
		}

		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var output = new(jobs.Jobs)
			var client = new(UserModel)
			result = jq.db.Where("id = ?", element.ClientID).First(&client)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client, 404")
			}
			output.ID = element.ID
			output.WorkerName = worker.Nama

			output.ClientName = client.Nama
			output.Foto = client.Foto

			var skill = new(SkillModel)
			result = jq.db.Where("id = ?", element.Category).First(&skill)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client")
			}
			output.Category = skill.NamaSkill
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price

			output.Status = element.Status

			*outputs = append(*outputs, *output)
		}
		return *outputs, int(totalCount), nil
	case "client":
		var client = new(UserModel)
		result := jq.db.Where("id = ?", userID).First(&client)
		if result.Error != nil {
			return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client, 404")
		}

		if role != client.Role {
			return nil, 0, errors.New("salah token")
		}

		if err := jq.db.
			Where("client_id = ? AND status = ?", userID, status).
			Find(&prePagination).
			Count(&totalCount).Error; err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil, 0, nil
			}
			return nil, 0, errors.New("server error")
		}
		if err := jq.db.
			Where("client_id = ? AND status = ?", userID, status).
			Order("updated_at desc").
			Offset(offset).
			Limit(pagesize).
			Find(&proses).
			Error; err != nil {
			return nil, 0, errors.New("server error")
		}
		if len(*proses) == 0 {
			return nil, int(totalCount), nil
		}

		var outputs = new([]jobs.Jobs)
		for _, element := range *proses {
			var worker = new(UserModel)
			result = jq.db.Where("id = ?", element.WorkerID).First(&worker)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan worker, 404")
			}
			var output = new(jobs.Jobs)
			output.WorkerName = worker.Nama
			output.ID = element.ID
			output.ClientName = client.Nama
			output.Foto = worker.Foto

			var skill = new(SkillModel)
			result = jq.db.Where("id = ?", element.Category).First(&skill)
			if result.Error != nil {
				return []jobs.Jobs{}, 0, errors.New("tidak ditemukan client")
			}
			output.Category = skill.NamaSkill
			output.StartDate = element.StartDate
			output.EndDate = element.EndDate
			output.Price = element.Price

			output.Status = element.Status

			*outputs = append(*outputs, *output)
		}
		return *outputs, int(totalCount), nil
	default:
		return nil, 0, errors.New("kesalahan pada role")
	}
}

func (jq *jobQuery) GetJob(jobID uint, role string) (jobs.Jobs, error) {
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
	result = jq.db.Where("id = ?", proses.WorkerID).First(&worker)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan worker, 404")
	}

	if role == "client" {
		output.Foto = worker.Foto
		output.NoHp = worker.NoHp
	} else if role == "worker" {
		output.Foto = client.Foto
		output.NoHp = client.NoHp
	}
	output.ID = proses.ID

	var skill = new(SkillModel)
	result = jq.db.Where("id = ?", proses.Category).First(&skill)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan client")
	}

	var statusPay = new(Transaction)
	result = jq.db.Where("job_id = ?", proses.ID).First(&statusPay)
	if result.Error != nil {
		output.StatusPayment = ""
	}

	output.StatusPayment = statusPay.Status

	output.Category = skill.NamaSkill

	output.WorkerName = worker.Nama

	output.ClientName = client.Nama

	output.StartDate = proses.StartDate
	output.EndDate = proses.EndDate
	output.Address = proses.Address
	output.Price = proses.Price
	output.Deskripsi = proses.Deskripsi
	output.Note = proses.NoteNego
	output.Status = proses.Status

	return *output, nil

}
func (jq *jobQuery) UpdateJob(update jobs.Jobs) (jobs.Jobs, error) {
	var proses = new(JobModel)
	result := jq.db.Where("id = ?", update.ID).First(&proses)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan jobs")
	}
	if proses.Status == "accepted" {
		if update.Status != "finished" {
			return jobs.Jobs{}, errors.New("jobs tidak boleh diubah, 403")
		}

	}

	if update.Status == "accepted" {
		update.Price = 0
	}

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
	if update.Note != "" {
		proses.NoteNego = update.Note
	}
	if update.Status != "" {
		proses.Status = update.Status
	}
	result = jq.db.Save(&proses)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("eror saat menyimpan data, 500")
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
	if update.Role == "client" {
		output.Foto = worker.Foto
		output.NoHp = worker.NoHp
	} else if update.Role == "worker" {
		output.Foto = client.Foto
		output.NoHp = client.NoHp
	}

	output.WorkerName = worker.Nama

	output.ClientName = client.Nama

	var skill = new(SkillModel)
	result = jq.db.Where("id = ?", proses.Category).First(&skill)
	if result.Error != nil {
		return jobs.Jobs{}, errors.New("tidak ditemukan client")
	}
	output.Category = skill.NamaSkill
	output.StartDate = proses.StartDate
	output.EndDate = proses.EndDate
	output.Price = proses.Price
	output.Deskripsi = proses.Deskripsi
	output.Status = proses.Status
	output.Address = proses.Address
	output.Note = proses.NoteNego

	return *output, nil

}
