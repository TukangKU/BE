package services

import (
	"errors"
	"tukangku/features/notifications"
)

type notifService struct {
	repo notifications.Repository
}

func New(r notifications.Repository) notifications.Service {
	return &notifService{
		repo: r,
	}
}

func (ns *notifService) GetNotifs(uid uint) ([]notifications.Notif, error) {
	// cek role

	// bikin di repo
	result, err := ns.repo.GetNotifs(uid)

	if err != nil {

		return nil, errors.New("terjadi kesalahan pada sistem")
	}
	// fmt.Println(result, "service")

	return result, nil
}
