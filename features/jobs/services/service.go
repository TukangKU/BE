package services

import (
	"errors"
	"fmt"
	"strings"
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
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return jobs.Jobs{}, errors.New("not found")
		}

		return jobs.Jobs{}, errors.New("terjadi kesalahan pada sistem")
	}
	// fmt.Println(result, "service")
	return result, nil
}

func (js *jobsService) GetJobs(id uint, status string, role string, page int, pagesize int) ([]jobs.Jobs, int, error) {
	if status == "" {
		// code jika tidak pake query
		result, count, err := js.repo.GetJobs(id, role, page, pagesize)
		if err != nil {
			// eror handling
			return nil, 0, err
		}
		return result, count, nil
	}

	result, count, err := js.repo.GetJobsByStatus(id, status, role, page, pagesize)
	if err != nil {
		// eror handling
		return nil, 0, err
	}
	return result, count, nil
}

func (js *jobsService) GetJob(jobID uint, role string) (jobs.Jobs, error) {
	if role != "worker" {
		if role != "client" {
			return jobs.Jobs{}, errors.New("role tidak dikenali")
		}
	}
	result, err := js.repo.GetJob(jobID, role)
	if err != nil {
		// eror handling
		return jobs.Jobs{}, err
	}

	fmt.Println(result, "servis")
	return result, nil
}

func (js *jobsService) UpdateJob(update jobs.Jobs) (jobs.Jobs, error) {
	// cek role
	// if update.Role == "client" {
	// 	update.Price = 0
	// 	update.Status = ""
	// }

	result, err := js.repo.UpdateJob(update)
	if err != nil {
		// eror handling
		return jobs.Jobs{}, err
	}
	return result, nil
}
