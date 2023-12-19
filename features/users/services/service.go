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
	if newUser.Email == "" || newUser.Password == "" {
		return users.Users{}, errors.New("incorrect inpu data")
	}

	ePassword, err := ur.h.HashPassword(newUser.Password)

	if err != nil {
		return users.Users{}, errors.New("terdapat masalah saat memproses data")
	}

	newUser.Password = ePassword
	result, err := ur.repo.Register(newUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return users.Users{}, errors.New("data telah terdaftar pada sistem")
		}
		return users.Users{}, errors.New("terjadi kesalahan pada sistem")
	}

	return result, nil
}

func (ul *userService) Login(email string, password string) (users.Users, error) {
	if email == "" || password == "" {
		return users.Users{}, errors.New("incorrect input data")
	}

	result, err := ul.repo.Login(email)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return users.Users{}, errors.New("data tidak ditemukan")
		}
		return users.Users{}, errors.New("terjadi kesalahan pada sistem")
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
