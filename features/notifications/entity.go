package notifications

import "github.com/labstack/echo/v4"

type Notif struct {
	ID        uint
	UserID    uint
	Message   string
	CreatedAt string
}

type Handler interface {
	GetNotifs() echo.HandlerFunc
}

type Service interface {
	GetNotifs(userid uint) ([]Notif, error)
}

type Repository interface {
	GetNotifs(userid uint) ([]Notif, error)
}
