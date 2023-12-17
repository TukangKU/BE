package users

import (
	"tukangku/features/skill"

	"github.com/labstack/echo/v4"
)

type Users struct {
	ID       uint           `json:"id"`
	Nama     string         `json:"nama"`
	UserName string         `json:"username"`
	Email    string         `json:"email"`
	Alamat   string         `json:"alamat"`
	NoHp     string         `json:"nohp"`
	Password string         `json:"password"`
	Foto     string         `json:"foto"`
	Skill    []skill.Skills `json:"skill"`
	Role     string         `json:"role"`
}

type Handler interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
}

type Service interface {
	Register(newUser Users) (Users, error)
	Login(email string, password string) (Users, error)
	UpdateUser(idUser uint, updateWorker Users) (Users, error)
}

type Repository interface {
	Register(newUser Users) (Users, error)
	Login(email string) (Users, error)
	UpdateUser(idUser uint, updateWorker Users) (Users, error)
}
