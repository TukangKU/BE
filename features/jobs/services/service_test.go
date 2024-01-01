package services_test

import (
	"errors"
	"testing"
	"tukangku/features/jobs"
	"tukangku/features/jobs/mocks"
	"tukangku/features/jobs/services"

	"github.com/stretchr/testify/assert"
)

func TestCreateJob(t *testing.T) {
	repo := mocks.NewRepository(t)
	srv := services.New(repo)

	t.Run("CheckClient", func(t *testing.T) {
		createJobFalse := jobs.Jobs{
			ID:       uint(1),
			ClientID: uint(1),
			Role:     "worker",
			SkillID:  uint(3),
		}

		res, err := srv.Create(createJobFalse)

		assert.Error(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
	})

	t.Run("Repository Error", func(t *testing.T) {
		createJob := jobs.Jobs{
			WorkerID:  uint(1),
			ClientID:  uint(1),
			Role:      "client",
			SkillID:   uint(3),
			StartDate: "20-20-2020",
			EndDate:   "21-21-2121",
		}

		repo.On("Create", createJob).Return(jobs.Jobs{}, errors.New("repository error")).Once()

		res, err := srv.Create(createJob)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
	})

	t.Run("Not Found Error", func(t *testing.T) {
		createJob := jobs.Jobs{
			WorkerID:  uint(1),
			ClientID:  uint(1),
			Role:      "client",
			SkillID:   uint(3),
			StartDate: "20-20-2020",
			EndDate:   "21-21-2121",
		}

		repo.On("Create", createJob).Return(jobs.Jobs{}, errors.New("tidak ditemukan")).Once()

		res, err := srv.Create(createJob)

		assert.Error(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
		assert.Equal(t, "not found", err.Error())
	})

	t.Run("Success Case", func(t *testing.T) {

		createJob := jobs.Jobs{
			WorkerID:  uint(1),
			ClientID:  uint(1),
			Role:      "client",
			SkillID:   uint(3),
			StartDate: "20-20-2020",
			EndDate:   "21-21-2121",
		}

		repo.On("Create", createJob).Return(jobs.Jobs{}, nil).Once()

		res, err := srv.Create(createJob)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
	})
}

func TestGetJobs(t *testing.T) {
	repo := mocks.NewRepository(t)
	srv := services.New(repo)

	t.Run("Status Error", func(t *testing.T) {
		id := uint(1)
		status := ""
		role := ""
		page := 1
		pageSize := 1
		count := 0

		repo.On("GetJobs", id, role, page, pageSize).Return(nil, count, errors.New("not found")).Once()

		res, x, err := srv.GetJobs(id, role, status, page, pageSize)

		repo.AssertExpectations(t)
		assert.Equal(t, err, errors.New("not found"))
		assert.Error(t, err)
		assert.Equal(t, x, count)
		assert.Nil(t, res)
	})

	t.Run("Status Error", func(t *testing.T) {
		succesReturnN := []jobs.Jobs{
			{
				ID:         88,
				WorkerName: "Paijo",
				ClientName: "Malik",
				Category:   "Service AC",
				Foto:       "malik.png",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      500000,
				Status:     "",
			},
			{
				ID:         77,
				WorkerName: "Paijo",
				ClientName: "Hafidz",
				Category:   "Plumber",
				Foto:       "hafidz.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      0,
				Status:     "",
			},
		}

		id := uint(1)
		status := ""
		role := ""
		page := 1
		pageSize := 1
		count := 1

		repo.On("GetJobs", id, role, page, pageSize).Return(succesReturnN, count, nil).Once()

		res, x, err := srv.GetJobs(id, role, status, page, pageSize)

		repo.AssertExpectations(t)
		assert.Equal(t, res, succesReturnN)
		assert.Equal(t, x, count)
		assert.Nil(t, err)
	})

	t.Run("Repo Error", func(t *testing.T) {
		id := uint(1)
		status := "pending"
		role := "worker"
		page := 1
		pageSize := 1
		count := 0

		repo.On("GetJobsByStatus", id, status, role, page, pageSize).Return(nil, count, errors.New("not found")).Once()

		res, x, err := srv.GetJobs(id, status, role, page, pageSize)

		repo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, x, count)
		assert.Nil(t, res)
	})

	t.Run("Repository Error", func(t *testing.T) {
		succesReturnN := []jobs.Jobs{
			{
				ID:         88,
				WorkerName: "Paijo",
				ClientName: "Malik",
				Category:   "Service AC",
				Foto:       "malik.png",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      500000,
				Status:     "",
			},
			{
				ID:         77,
				WorkerName: "Paijo",
				ClientName: "Hafidz",
				Category:   "Plumber",
				Foto:       "hafidz.jpg",
				StartDate:  "25-12-2023",
				EndDate:    "25-12-2023",
				Price:      0,
				Status:     "",
			},
		}

		id := uint(1)
		status := "pending"
		role := "worker"
		page := 1
		pageSize := 1
		count := 0

		repo.On("GetJobsByStatus", id, status, role, page, pageSize).Return(succesReturnN, count, nil).Once()

		res, x, err := srv.GetJobs(id, status, role, page, pageSize)

		repo.AssertExpectations(t)
		assert.Equal(t, res, succesReturnN)
		assert.Equal(t, x, count)
		assert.Nil(t, err)
	})
}

func TestGetJob(t *testing.T) {
	repo := mocks.NewRepository(t)
	srv := services.New(repo)

	t.Run("Wrong Role", func(t *testing.T) {
		id := uint(1)
		role := ""

		res, err := srv.GetJob(id, role)

		assert.Error(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
	})

	t.Run("Repository error", func(t *testing.T) {
		id := uint(1)
		role := "worker"

		repo.On("GetJob", id, role).Return(jobs.Jobs{}, errors.New("repository error")).Once()

		res, err := srv.GetJob(id, role)

		repo.AssertExpectations(t)
		
		assert.Error(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
	})

	t.Run("Success Case", func(t *testing.T) {
		id := uint(1)
		role := "worker"

		repo.On("GetJob", id, role).Return(jobs.Jobs{}, nil).Once()

		res, err := srv.GetJob(id, role)

		repo.AssertExpectations(t)
		
		assert.Nil(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
	})

}


func TestUpdateJob(t *testing.T) {
	repo := mocks.NewRepository(t)
	srv := services.New(repo)

	t.Run("Repository Error", func(t *testing.T) {
		inputUpdate := jobs.Jobs{
			Price: 20000,
			Deskripsi: "some string",
			Note: "some string",
		}

		repo.On("UpdateJob", inputUpdate).Return(jobs.Jobs{}, errors.New("repository error")).Once()

		res, err := srv.UpdateJob(inputUpdate)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
		
	})


	t.Run("Success Case", func(t *testing.T) {
		inputUpdate := jobs.Jobs{
			Price: 20000,
			Deskripsi: "some string",
			Note: "some string",
		}

		repo.On("UpdateJob", inputUpdate).Return(jobs.Jobs{}, nil).Once()

		res, err := srv.UpdateJob(inputUpdate)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, jobs.Jobs{}, res)
		
	})
}