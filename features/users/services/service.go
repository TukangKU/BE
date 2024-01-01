package services

import (
	"errors"
	"strings"
	"tukangku/features/users"
	"tukangku/helper/enkrip"
)

type userService struct {
	repo users.Repository
	h    enkrip.HashInterface
}

func New(r users.Repository, h enkrip.HashInterface) users.Service {
	return &userService{
		repo: r,
		h:    h,
	}
}

func (ur *userService) Register(newUser users.Users) (users.Users, error) {
	if newUser.Email == ""  {
		return users.Users{}, errors.New("email harus di isi")
	}

	if newUser.UserName == ""  {
		return users.Users{}, errors.New("username harus di isi")
	}

	if newUser.Password == ""  {
		return users.Users{}, errors.New("password harus di isi")
	}

	if newUser.Role == ""  {
		return users.Users{}, errors.New("role harus di isi")
	}

	ePassword, err := ur.h.HashPassword(newUser.Password)

	if err != nil {
		return users.Users{}, errors.New("terdapat masalah saat memproses data")
	}

	newUser.Password = ePassword
	result, err := ur.repo.Register(newUser)

	if err != nil {
		return users.Users{}, err
	}

	return result, nil
}

func (ul *userService) Login(email string, password string) (users.Users, error) {
	if email == "" {
		return users.Users{}, errors.New("email harus di isi")
	}

	if password == "" {
		return users.Users{}, errors.New("password harus di isi")
	}

	result, err := ul.repo.Login(email)

	if err != nil {
		return users.Users{}, err
	}

	err = ul.h.Compare(result.Password, password)

	if err != nil {
		return users.Users{}, errors.New("password salah")
	}

	return result, nil
}

func (us *userService) UpdateUser(idUser uint, updateWorker users.Users) (users.Users, error) {
	result, err := us.repo.UpdateUser(idUser, updateWorker)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return users.Users{}, errors.New("failed to update user")
		}
		return users.Users{}, errors.New("failed to update user")
	}
	return result, nil
}

func (gu *userService) GetUserByID(idUser uint) (users.Users, error) {
	result, err := gu.repo.GetUserByID(idUser)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return users.Users{}, errors.New("user not found")
		}
		return users.Users{}, errors.New("error retrieving User by ID")
	}

	return result, nil
}

func (gu *userService) GetUserBySKill(idSkill uint, page, pageSize int) ([]users.Users, int, error) {
	result, totalCount, err := gu.repo.GetUserBySKill(idSkill, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return result, totalCount, nil
}

