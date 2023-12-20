package services

import (
	"errors"
	"tukangku/features/jobs"
)

type jobsService struct {
	repo jobs.Repository
}

func New(r jobs.Repository) jobs.Service {
	return &jobsService{
		repo: r,
	}
}

func (js *jobsService) Create(newJobs jobs.Jobs) (jobs.Jobs, error) {
	// cek role
	if newJobs.Role != "client" {

		return jobs.Jobs{}, errors.New("anda bukan client")

	}
	// bikin di repo
	result, err := js.repo.Create(newJobs)

	if err != nil {

		return jobs.Jobs{}, errors.New("terjadi kesalahan pada sistem")
	}
	// fmt.Println(result, "service")
	return result, nil
}

func (js *jobsService) GetJobs(id uint, status string, role string) ([]jobs.Jobs, error) {
	if status == "" {
		// code jika tidak pake query
		result, err := js.repo.GetJobs(id, role)
		if err != nil {
			// eror handling
			return nil, err
		}
		return result, nil
	}

	result, err := js.repo.GetJobsByStatus(id, status, role)
	if err != nil {
		// eror handling
		return nil, err
	}
	return result, nil
}
