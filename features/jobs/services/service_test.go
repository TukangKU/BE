package services_test

import (
	"errors"
	"testing"
	"tukangku/features/jobs"
	"tukangku/features/jobs/services"
	"tukangku/features/mocks"

	"github.com/stretchr/testify/assert"
)

var workerID uint = 1

var clientID uint = 2

var skillID uint = 3
var startDate, endDate, deskripsi, alamat = "2023-12-25", "2023-12-25", "Mas, tolong benerin sambungan pipa ke wastafel", "Jl.Setiabudi nomor 3"

func TestCreate(t *testing.T) {
	repo := mocks.NewRepository(t)
	msrv := mocks.NewService(t)
	srv := services.New(repo)
	var newJob = jobs.Jobs{
		WorkerID:  workerID,
		ClientID:  clientID,
		StartDate: startDate,
		EndDate:   endDate,
		SkillID:   skillID,
		Deskripsi: deskripsi,
		Address:   alamat,
		Role:      "client",
	}

	// no worker
	var noWorker = jobs.Jobs{

		ClientID:  clientID,
		StartDate: startDate,
		EndDate:   endDate,
		SkillID:   skillID,
		Deskripsi: deskripsi,
		Address:   alamat,
		Role:      "client",
	}
	// no client
	var noCLient = jobs.Jobs{
		WorkerID: workerID,

		StartDate: startDate,
		EndDate:   endDate,
		SkillID:   skillID,
		Deskripsi: deskripsi,
		Address:   alamat,
		Role:      "client",
	}
	// no category
	var noSkillID = jobs.Jobs{
		WorkerID:  workerID,
		ClientID:  clientID,
		StartDate: startDate,
		EndDate:   endDate,

		Deskripsi: deskripsi,
		Address:   alamat,
		Role:      "client",
	}
	// no start
	var noStartDate = jobs.Jobs{
		WorkerID: workerID,
		ClientID: clientID,

		EndDate:   endDate,
		SkillID:   skillID,
		Deskripsi: deskripsi,
		Address:   alamat,
		Role:      "client",
	}
	// no end
	var noEndDate = jobs.Jobs{
		WorkerID:  workerID,
		ClientID:  clientID,
		StartDate: startDate,

		SkillID:   skillID,
		Deskripsi: deskripsi,
		Address:   alamat,
		Role:      "client",
	}
	// wrong role
	var wrongRole = jobs.Jobs{
		WorkerID:  workerID,
		ClientID:  clientID,
		StartDate: startDate,
		EndDate:   endDate,
		SkillID:   skillID,
		Deskripsi: deskripsi,
		Address:   alamat,
		Role:      "worker",
	}
	//
	//

	var result = jobs.Jobs{
		ID:         1,
		Foto:       "worker.jpg",
		WorkerName: "Paijo",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     0,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "pending",
		Address:   "Jl.Setiabudi nomor 3",
	}
	t.Run("Success Case", func(t *testing.T) {
		repo.On("Create", newJob).Return(result, nil).Once()
		proses, err := srv.Create(newJob)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, result, proses)
	})
	t.Run("Case 1", func(t *testing.T) {
		repo.On("Create", noWorker).Return(jobs.Jobs{}, errors.New("please input worker_id")).Once()
		data, err := srv.Create(noWorker)
		assert.Error(t, err)
		assert.Equal(t, "please input worker_id", err.Error())
		assert.Equal(t, data, jobs.Jobs{})
		repo.AssertExpectations(t)
	})
	t.Run("Case 2", func(t *testing.T) {
		repo.On("Create", noCLient).Return(jobs.Jobs{}, errors.New("please input client_id")).Once()
		data, err := srv.Create(noCLient)
		assert.Error(t, err)
		assert.Equal(t, "please input client_id", err.Error())
		assert.Equal(t, data, jobs.Jobs{})
		repo.AssertExpectations(t)
	})
	t.Run("Case 3", func(t *testing.T) {
		repo.On("Create", noSkillID).Return(jobs.Jobs{}, errors.New("please input skill_id")).Once()
		data, err := srv.Create(noSkillID)
		assert.Error(t, err)
		assert.Equal(t, "please input skill_id", err.Error())
		assert.Equal(t, data, jobs.Jobs{})
		repo.AssertExpectations(t)
	})
	t.Run("Case 4", func(t *testing.T) {
		repo.On("Create", noStartDate).Return(jobs.Jobs{}, errors.New("please input start_date")).Once()
		data, err := srv.Create(noStartDate)
		assert.Error(t, err)
		assert.Equal(t, "please input start_date", err.Error())
		assert.Equal(t, data, jobs.Jobs{})
		repo.AssertExpectations(t)
	})
	t.Run("Case 5", func(t *testing.T) {
		repo.On("Create", noEndDate).Return(jobs.Jobs{}, errors.New("please input end_date")).Once()
		data, err := srv.Create(noEndDate)
		assert.Error(t, err)
		assert.Equal(t, "please input end_date", err.Error())
		assert.Equal(t, data, jobs.Jobs{})
		repo.AssertExpectations(t)
	})
	t.Run("Case 6", func(t *testing.T) {
		msrv.On("Create", wrongRole).Return(jobs.Jobs{}, errors.New("anda bukan client"))
		a, err := msrv.Create(wrongRole)
		assert.Error(t, err)
		assert.Equal(t, "anda bukan client", err.Error())
		assert.Equal(t, a, jobs.Jobs{})
		repo.AssertExpectations(t)
	})
}

func TestGetJob(t *testing.T) {
	repo := mocks.NewRepository(t)
	msrv := mocks.NewService(t)
	srv := services.New(repo)
	var jobID uint = 1
	var jobIDFalse uint = 2
	var resultA = jobs.Jobs{
		ID:         1,
		Foto:       "worker.jpg",
		WorkerName: "Paijo",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     0,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "pending",
		Address:   "Jl.Setiabudi nomor 3",
	}
	var resultB = jobs.Jobs{
		ID:         1,
		Foto:       "client.jpg",
		WorkerName: "Paijo",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     0,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "pending",
		Address:   "Jl.Setiabudi nomor 3",
	}
	roleA, roleB := "client", "worker"
	t.Run("Success Case 1", func(t *testing.T) {
		repo.On("GetJob", jobID, roleA).Return(resultA, nil).Once()
		proses, err := srv.GetJob(jobID, roleA)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, resultA, proses)
	})
	t.Run("Success Case 2", func(t *testing.T) {
		repo.On("GetJob", jobID, roleB).Return(resultB, nil).Once()
		proses, err := srv.GetJob(jobID, roleB)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, resultB, proses)
	})
	t.Run("Failed Case", func(t *testing.T) {
		repo.On("GetJob", jobIDFalse, roleA).Return(jobs.Jobs{}, errors.New("not found")).Once()
		proses, err := srv.GetJob(jobIDFalse, roleA)
		repo.AssertExpectations(t)
		assert.Equal(t, proses, jobs.Jobs{})

		assert.Equal(t, errors.New("not found"), err)
	})
	t.Run("Failed Case 2", func(t *testing.T) {
		repo.On("GetJob", jobIDFalse, roleB).Return(jobs.Jobs{}, errors.New("not found")).Once()
		proses, err := srv.GetJob(jobIDFalse, roleB)
		repo.AssertExpectations(t)
		assert.Equal(t, proses, jobs.Jobs{})

		assert.Equal(t, errors.New("not found"), err)
	})
	t.Run("Failed Case 3", func(t *testing.T) {
		msrv.On("GetJob", jobIDFalse, "roleB").Return(jobs.Jobs{}, errors.New("role tidak dikenali")).Once()
		proses, err := msrv.GetJob(jobIDFalse, "roleB")
		repo.AssertExpectations(t)
		assert.Equal(t, proses, jobs.Jobs{})

		assert.Equal(t, errors.New("role tidak dikenali"), err)
	})
}
func TestUpdateJob(t *testing.T) {
	repo := mocks.NewRepository(t)
	// msrv := mocks.NewService(t)
	srv := services.New(repo)
	var jobID uint = 1
	var reqBodyW = jobs.Jobs{
		ID:   1,
		Role: "worker",

		WorkerID: workerID,
		Price:    300000,

		Status: "negotiation",
		Note:   "baik mas, harganya 300000",
	}
	// var jobIDFalse uint = 2
	// var resultA = jobs.Jobs{
	// 	ID:         1,
	// 	Foto:       "worker.jpg",
	// 	WorkerName: "Paijo",
	// 	Category:   "Plumber",

	// 	StartDate: "2023-12-25",
	// 	EndDate:   "2023-12-25",
	// 	Price:     0,
	// 	Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
	// 	Status:    "pending",
	// 	Address:   "Jl.Setiabudi nomor 3",
	// }
	var resultW = jobs.Jobs{
		ID:         jobID,
		Foto:       "client.jpg",
		WorkerName: "Paijo",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     300000,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "negotiation",
		Note:      "baik mas, harganya 300000",
		Address:   "Jl.Setiabudi nomor 3",
	}
	t.Run("Success Case W", func(t *testing.T) {
		repo.On("UpdateJob", reqBodyW).Return(resultW, nil).Once()
		proses, err := srv.UpdateJob(reqBodyW)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, resultW, proses)
	})
}
