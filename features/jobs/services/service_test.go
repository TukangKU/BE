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
	var reqBodyC = jobs.Jobs{
		ID:   1,
		Role: "client",

		ClientID: clientID,
		Price:    250000,

		Status: "negotiation",
		Note:   "250 aja ya mas?",
	}
	var reqBodyAcc = jobs.Jobs{
		ID:   1,
		Role: "worker",

		WorkerID: workerID,
		Price:    250000,

		Status: "accepted",
		Note:   "250 aja ya mas?",
	}
	var reqBodyAccClient = jobs.Jobs{
		ID:   1,
		Role: "client",

		WorkerID: workerID,
		Price:    0,

		Status: "accepted",
	}
	var reqBodyNoID = jobs.Jobs{

		Role: "client",

		WorkerID: workerID,
		Price:    250000,

		Status: "negotiation",
		Note:   "250 aja ya mas?",
	}

	var reqBodyNoRole = jobs.Jobs{
		ID:   1,
		Role: "",

		WorkerID: workerID,
		Price:    250000,

		Status: "negotiation",
		Note:   "250 aja ya mas?",
	}
	// var jobIDFalse uint = 2
	var resultClient = jobs.Jobs{
		ID:         1,
		Foto:       "worker.jpg",
		WorkerName: "Paijo",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     250000,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "negotiation",
		Note:      "250 aja ya mas?",
		Address:   "Jl.Setiabudi nomor 3",
	}
	var resultClientNoPrice = jobs.Jobs{
		ID:         1,
		Foto:       "worker.jpg",
		WorkerName: "Paijo",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     300000,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "accepted",
		Note:      "oke mas",
		Address:   "Jl.Setiabudi nomor 3",
	}
	var resultWorker = jobs.Jobs{
		ID:         jobID,
		Foto:       "client.jpg",
		ClientName: "Johan",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     300000,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "negotiation",
		Note:      "baik mas, harganya 300000",
		Address:   "Jl.Setiabudi nomor 3",
	}
	var resultAcc = jobs.Jobs{
		ID:         jobID,
		Foto:       "client.jpg",
		WorkerName: "Paijo",
		Category:   "Plumber",

		StartDate: "2023-12-25",
		EndDate:   "2023-12-25",
		Price:     300000,
		Deskripsi: "Mas, tolong benerin sambungan pipa ke wastafel",
		Status:    "accepted",
		Note:      "250 aja ya mas?",
		Address:   "Jl.Setiabudi nomor 3",
	}
	t.Run("Success Case Worker", func(t *testing.T) {
		repo.On("UpdateJob", reqBodyW).Return(resultWorker, nil).Once()
		proses, err := srv.UpdateJob(reqBodyW)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, resultWorker, proses)
	})
	t.Run("Success Case Client", func(t *testing.T) {
		repo.On("UpdateJob", reqBodyC).Return(resultClient, nil).Once()
		proses, err := srv.UpdateJob(reqBodyC)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, resultClient, proses)
	})
	t.Run("Success Case Client No Price", func(t *testing.T) {
		repo.On("UpdateJob", reqBodyAccClient).Return(resultClientNoPrice, nil).Once()
		proses, err := srv.UpdateJob(reqBodyAccClient)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, resultClientNoPrice, proses)
	})
	t.Run("Success Case Acc", func(t *testing.T) {
		repo.On("UpdateJob", reqBodyAcc).Return(resultAcc, nil).Once()
		proses, err := srv.UpdateJob(reqBodyAcc)
		repo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, resultAcc, proses)
	})

	t.Run("Failed Case no ID", func(t *testing.T) {
		repo.On("UpdateJob", reqBodyNoID).Return(jobs.Jobs{}, errors.New("tidak ditemukan jobs"))
		result, err := srv.UpdateJob(reqBodyNoID)
		repo.AssertExpectations(t)
		assert.Equal(t, result, jobs.Jobs{})
		assert.Equal(t, err, errors.New("tidak ditemukan jobs"))
	})

	t.Run("Failed Case no Role", func(t *testing.T) {
		repo.On("UpdateJob", reqBodyNoRole).Return(jobs.Jobs{}, errors.New("masukkan role, 403"))
		result, err := srv.UpdateJob(reqBodyNoRole)
		repo.AssertExpectations(t)
		assert.Equal(t, result, jobs.Jobs{})
		assert.Equal(t, err, errors.New("masukkan role, 403"))
	})

	t.Run("Failed Case accepted", func(t *testing.T) {

		repo.On("UpdateJob", reqBodyC).Return(jobs.Jobs{}, errors.New("jobs tidak boleh diubah, 403"))
		result, err := srv.UpdateJob(reqBodyC)
		repo.AssertExpectations(t)
		assert.Equal(t, result, jobs.Jobs{})
		assert.Equal(t, err, errors.New("jobs tidak boleh diubah, 403"))
	})
}

func TestGetJobs(t *testing.T) {
	repo := mocks.NewRepository(t)
	// msrv := mocks.NewService(t)
	srv := services.New(repo)
	// var jobID uint = 1

	t.Run("success no query", func(t *testing.T) {
		succesReturnN := []jobs.Jobs{
			{
				ID:         32,
				WorkerName: "Julian",
				ClientName: "Johan",
				Category:   "Service AC",
				Foto:       "julian.png",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      500000,
				Status:     "accepted",
			},
			{
				ID:         31,
				WorkerName: "Paijo",
				ClientName: "Johan",
				Category:   "Plumber",
				Foto:       "paijo.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      0,
				Status:     "pending",
			},
			{
				ID:         30,
				WorkerName: "Lazuardi",
				ClientName: "Johan",
				Category:   "CCTV",
				Foto:       "lazu.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      5000000,
				Status:     "finished",
			},
		}
		count := 3
		page := 1
		pagesize := 10
		repo.On("GetJobs", clientID, "client", page, pagesize).Return(succesReturnN, count, nil)
		result, x, err := repo.GetJobs(clientID, "client", page, pagesize)

		repo.AssertExpectations(t)
		assert.Equal(t, result, succesReturnN)
		assert.Equal(t, x, count)
		assert.Nil(t, err)

	})
	t.Run("succes with param", func(t *testing.T) {
		succesReturn := []jobs.Jobs{

			{
				ID:         30,
				WorkerName: "Lazuardi",
				ClientName: "Johan",
				Category:   "CCTV",
				Foto:       "lazu.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      5000000,
				Status:     "finished",
			},
		}
		succesReturnN := []jobs.Jobs{
			{
				ID:         32,
				WorkerName: "Julian",
				ClientName: "Johan",
				Category:   "Service AC",
				Foto:       "julian.png",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      500000,
				Status:     "accepted",
			},
			{
				ID:         31,
				WorkerName: "Paijo",
				ClientName: "Johan",
				Category:   "Plumber",
				Foto:       "paijo.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      0,
				Status:     "pending",
			},
			{
				ID:         30,
				WorkerName: "Lazuardi",
				ClientName: "Johan",
				Category:   "CCTV",
				Foto:       "lazu.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      5000000,
				Status:     "finished",
			},
		}
		count := 3
		countB := 1
		page := 1
		pagesize := 10
		status := "finished"

		// repo.On("GetJobs", clientID, "client", status, page, pagesize).Return(succesReturn, countB, nil)
		repo.On("GetJobsByStatus", clientID, "client", status, page, pagesize).Return(succesReturn, countB, nil)
		repo.Mock.On("GetJobsByStatus", clientID, "client", status, page, pagesize).Return(succesReturn, countB, nil)
		repo.Mock.On("GetJobs", clientID, "client", page, pagesize).Return(succesReturnN, count, nil)
		result, x, err := srv.GetJobs(clientID, "client", status, page, pagesize)
		hasil, y, arr := repo.GetJobs(clientID, "client", page, pagesize)
		aa, bb, cc := repo.GetJobsByStatus(clientID, "client", status, page, pagesize)
		repo.AssertExpectations(t)
		assert.Equal(t, result, succesReturn)
		assert.Equal(t, hasil, succesReturnN)

		assert.Equal(t, x, countB)
		assert.Equal(t, y, count)

		assert.Nil(t, err)
		assert.Nil(t, arr)

		assert.Equal(t, aa, succesReturn)
		assert.Equal(t, bb, countB)
		assert.Nil(t, cc)
	})
	t.Run("false id", func(t *testing.T) {
		var falseID uint = 4
		page := 1
		pagesize := 10
		count := 0
		repo.On("GetJobs", falseID, "client", page, pagesize).Return(nil, count, errors.New("not found"))
		result, x, err := repo.GetJobs(falseID, "client", page, pagesize)

		repo.AssertExpectations(t)
		assert.Equal(t, err, errors.New("not found"))
		assert.Equal(t, x, count)
		assert.Nil(t, result)
	})
	t.Run("false role", func(t *testing.T) {

		page := 1
		pagesize := 10
		count := 0
		repo.On("GetJobs", clientID, "worker", page, pagesize).Return(nil, count, errors.New("tidak ditemukan worker, 404"))
		result, x, err := repo.GetJobs(clientID, "worker", page, pagesize)

		repo.AssertExpectations(t)
		assert.Equal(t, err, errors.New("tidak ditemukan worker, 404"))
		assert.Equal(t, x, count)
		assert.Nil(t, result)
	})
	t.Run("no role", func(t *testing.T) {

		page := 1
		pagesize := 10
		count := 0
		repo.On("GetJobs", clientID, "", page, pagesize).Return(nil, count, errors.New("kesalahan pada role"))
		result, x, err := repo.GetJobs(clientID, "", page, pagesize)

		repo.AssertExpectations(t)
		assert.Equal(t, err, errors.New("kesalahan pada role"))
		assert.Equal(t, x, count)
		assert.Nil(t, result)
	})
	// })
	t.Run("succes with param worker", func(t *testing.T) {
		succesReturn := []jobs.Jobs{

			{
				ID:         30,
				WorkerName: "Lazuardi",
				ClientName: "Johan",
				Category:   "CCTV",
				Foto:       "lazu.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      0,
				Status:     "finished",
			},
		}
		succesReturnN := []jobs.Jobs{
			{
				ID:         32,
				WorkerName: "Julian",
				ClientName: "Johan",
				Category:   "Service AC",
				Foto:       "julian.png",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      500000,
				Status:     "accepted",
			},
			{
				ID:         31,
				WorkerName: "Paijo",
				ClientName: "Johan",
				Category:   "Plumber",
				Foto:       "paijo.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      0,
				Status:     "pending",
			},
			{
				ID:         30,
				WorkerName: "Lazuardi",
				ClientName: "Johan",
				Category:   "CCTV",
				Foto:       "lazu.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      5000000,
				Status:     "finished",
			},
		}
		count := 3
		countB := 1
		page := 1
		pagesize := 10
		status := "finished"

		// repo.On("GetJobs", clientID, "client", status, page, pagesize).Return(succesReturn, countB, nil)
		repo.On("GetJobsByStatus", workerID, "worker", status, page, pagesize).Return(succesReturn, countB, nil)
		repo.Mock.On("GetJobsByStatus", workerID, "worker", status, page, pagesize).Return(succesReturn, countB, nil)
		repo.Mock.On("GetJobs", workerID, "worker", page, pagesize).Return(succesReturnN, count, nil)
		result, x, err := srv.GetJobs(workerID, "worker", status, page, pagesize)
		hasil, y, arr := repo.GetJobs(workerID, "worker", page, pagesize)
		aa, bb, cc := repo.GetJobsByStatus(workerID, "worker", status, page, pagesize)
		repo.AssertExpectations(t)
		assert.Equal(t, result, succesReturn)
		assert.Equal(t, hasil, succesReturnN)

		assert.Equal(t, x, countB)
		assert.Equal(t, y, count)

		assert.Nil(t, err)
		assert.Nil(t, arr)

		assert.Equal(t, aa, succesReturn)
		assert.Equal(t, bb, countB)
		assert.Nil(t, cc)
	})
	t.Run("Success Zero result", func(t *testing.T) {

		countB := 0
		page := 1
		pagesize := 10
		status := "rejected"
		repo.On("GetJobsByStatus", workerID, "worker", status, page, pagesize).Return([]jobs.Jobs{}, countB, nil)
		repo.Mock.On("GetJobsByStatus", workerID, "worker", status, page, pagesize).Return([]jobs.Jobs{}, countB, nil)
		result, x, err := srv.GetJobs(workerID, "worker", status, page, pagesize)
		repo.AssertExpectations(t)
		assert.Equal(t, x, countB)
		assert.Equal(t, result, []jobs.Jobs{})

		assert.Nil(t, err)
	})
}
