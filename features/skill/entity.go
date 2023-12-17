package skill

import (
	"github.com/labstack/echo/v4"
)

type Skills struct {
	ID        uint   `json:"id"`
	NamaSkill string `json:"skill"`
}

type Handler interface {
	Add() echo.HandlerFunc
	Show() echo.HandlerFunc
}

type Service interface {
	AddSkill(newSkill Skills) (Skills, error)
	ShowSkill() ([]Skills, error)
}

type Repository interface {
	AddSkill(newSkill Skills) (Skills, error)
	ShowSkill() ([]Skills, error)
}
