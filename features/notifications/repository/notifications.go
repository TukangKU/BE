package model

import (
	"errors"
	"fmt"
	"tukangku/features/notifications"

	"gorm.io/gorm"
)

type NotifModel struct {
	gorm.Model
	UserID  uint `gorm:"not null"`
	Message string
}

type notifQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) notifications.Repository {
	return &notifQuery{
		db: db,
	}
}

func (nq *notifQuery) GetNotifs(id uint) ([]notifications.Notif, error) {
	var proses = new([]NotifModel)

	if err := nq.db.Where("user_id = ?", id).Order("created_at desc").Find(&proses).Error; err != nil {
		return nil, errors.New("server error")
	}
	if len(*proses) == 0 {
		return nil, nil
	}
	fmt.Println(proses, "repo")
	var result = new([]notifications.Notif)

	for _, element := range *proses {
		var newResult = new(notifications.Notif)
		newResult.ID = element.ID
		newResult.Message = element.Message
		newResult.CreatedAt = element.CreatedAt.String()
		*result = append(*result, *newResult)
	}
	// fmt.Println(result, "repo")
	return *result, nil
}
