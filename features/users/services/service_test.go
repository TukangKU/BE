package services_test

import (
	"errors"
	"testing"
	"tukangku/features/users"
	"tukangku/features/users/mocks"
	"tukangku/features/users/services"
	emock "tukangku/helper/enkrip/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	enc := emock.NewHashInterface(t)
	userService := services.New(repo, enc)

	t.Run("Success Case", func(t *testing.T) {
		inputData := users.Users{
			UserName: "bedul",
			Email:    "bedul@gmail.com",
			Password: "bedulganteng",
			Role:     "worker",
		}

		enc.On("HashPassword", inputData.Password).Return("some string", nil).Once()

		inputData.Password = "some string"
		repo.On("Register", inputData).Return(inputData, nil).Once()

		inputData.Password = "bedulganteng"
		res, err := userService.Register(inputData)

		assert.Nil(t, err)
		assert.Equal(t, "some string", res.Password)

		enc.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		inputData := users.Users{
			UserName: "bedul",
			Email:    "",
			Role:     "client",
			Password: "bedulganteng",
		}

		res, err := userService.Register(inputData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Invalid Username", func(t *testing.T) {
		inputData := users.Users{
			UserName: "",
			Email:    "bedul@gmail.com",
			Role:     "client",
			Password: "bedulganteng",
		}

		res, err := userService.Register(inputData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		inputData := users.Users{
			UserName: "bedul",
			Email:    "bedul@gmail.com",
			Role:     "client",
			Password: "",
		}

		res, err := userService.Register(inputData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Invalid Role", func(t *testing.T) {
		inputData := users.Users{
			UserName: "bedul",
			Email:    "bedul@gmail.com",
			Role:     "",
			Password: "bedulganteng",
		}

		res, err := userService.Register(inputData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Encryp Error", func(t *testing.T) {
		inputData := users.Users{
			Email:    "bedul@gmail.com",
			UserName: "bedul",
			Role:     "worker",
			Password: "bedul",
		}

		enc.On("HashPassword", inputData.Password).Return("", errors.New("terdapat masalah saat memproses data")).Once()

		res, err := userService.Register(inputData)

		assert.Error(t, err)
		assert.Equal(t, "terdapat masalah saat memproses data", err.Error())
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Repository Error", func(t *testing.T) {
		inputData := users.Users{
			Email:    "bedul@gmail.com",
			UserName: "bedul",
			Role:     "worker",
			Password: "bedul",
		}

		enc.On("HashPassword", inputData.Password).Return("bedul", nil).Once()

		repo.On("Register", inputData).Return(users.Users{}, errors.New("bedul")).Once()

		res, err := userService.Register(inputData)

		enc.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	enc := emock.NewHashInterface(t)
	userService := services.New(repo, enc)

	t.Run("Invalid Email", func(t *testing.T) {
		inputData := users.Users{
			Email:    "",
			Password: "bedulganteng",
		}

		res, err := userService.Login(inputData.Email, inputData.Password)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		inputData := users.Users{
			Email:    "bedul@gmail.com",
			Password: "",
		}

		res, err := userService.Login(inputData.Email, inputData.Password)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Repository Error", func(t *testing.T) {
		inputData := users.Users{
			Email:    "bedul@gmail.com",
			Password: "bedulganteng",
		}

		repo.On("Login", inputData.Email).Return(users.Users{}, errors.New("repository error")).Once()

		res, err := userService.Login(inputData.Email, inputData.Password)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		inputData := users.Users{
			Email:    "bedul@gmail.com",
			Password: "passwordsalah",
		}

		dataUser := users.Users{
			ID:       1,
			Nama:     "bedul",
			Email:    "bedul@gmail.com",
			Password: "bedul",
		}

		repo.On("Login", inputData.Email).Return(dataUser, nil).Once()
		enc.On("Compare", dataUser.Password, inputData.Password).Return(errors.New("password salah")).Once()
		res, err := userService.Login(inputData.Email, inputData.Password)

		enc.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)

	})

	t.Run("Success Case", func(t *testing.T) {
		inputData := users.Users{
			Email:    "bedul@gmail.com",
			Password: "bedul",
		}

		dataUser := users.Users{
			ID:       1,
			Nama:     "bedul",
			Email:    "bedul@gmail.com",
			Password: "bedul",
		}

		repo.On("Login", inputData.Email).Return(dataUser, nil).Once()
		enc.On("Compare", dataUser.Password, inputData.Password).Return(nil).Once()
		res, err := userService.Login(inputData.Email, inputData.Password)

		enc.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "bedul", res.Password)

	})
}

func TestUpdateUser(t *testing.T) {
	repo := mocks.NewRepository(t)
	userService := services.New(repo, nil)

	userID := uint(1)
	inputData := users.Users{
		UserName: "bedulganteng",
		Email:    "bedul@gmail.com",
	}
	updateData := users.Users{
		UserName: "bedulganteng",
		Email:    "bedul@gmail.com",
	}

	t.Run("Not Found", func(t *testing.T) {

		repo.On("UpdateUser", userID, inputData).Return(users.Users{}, errors.New("not found")).Once()

		res, err := userService.UpdateUser(userID, updateData)

		
		assert.Error(t, err)
		assert.Equal(t, users.Users{}, res)
		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - Other Error", func(t *testing.T) {
		repo.On("UpdateUser", userID, updateData).Return(users.Users{}, errors.New("repository error")).Once()

		result, err := userService.UpdateUser(userID, updateData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "failed to update user", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Success Case", func(t *testing.T) {
		
		repo.On("UpdateUser", userID, updateData).Return(users.Users{}, nil).Once()

		result, err := userService.UpdateUser(userID, updateData)

		assert.Nil(t, err)
		assert.Equal(t, users.Users{}, result)

		repo.AssertExpectations(t)
	})
}


func TestGetUserByID(t *testing.T) {
	var userID = uint(1)
	repo := mocks.NewRepository(t)
	enkrip := emock.NewHashInterface(t)
	userService := services.New(repo, enkrip)

	t.Run("Success Case", func(t *testing.T) {
		expectedResult := users.Users{
			ID:       userID,
			Nama:     "bedul",
			UserName: "bedul",
			Email:    "bedul@gmail.com",
			Foto:     "bedul.jpg",
		}

		repo.On("GetUserByID", userID).Return(expectedResult, nil).Once()

		result, err := userService.GetUserByID(userID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResult, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case", func(t *testing.T) {
		repo.On("GetUserByID", userID).Return(users.Users{
			ID:       0,
			UserName: "unknown",
			Email:    "unknown@example.com",
		}, errors.New("error retrieving User by ID")).Once()

		result, err := userService.GetUserByID(userID)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "error retrieving User by ID", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {

		repo.On("GetUserByID", userID).Return(users.Users{}, errors.New("not found")).Once()

		result, err := userService.GetUserByID(userID)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "user not found", err.Error())
	})

}


func TestGetUserBySkill(t *testing.T) {
	repo := mocks.NewRepository(t)
	service := services.New(repo, nil)

	idSkill := uint(1)
	page := 1
	pageSize := 10

	t.Run("Success Case", func(t *testing.T) {
		mockResult := []users.Users{
			{ID: 1, UserName: "user1", Email: "user1@example.com"},
			{ID: 2, UserName: "user2", Email: "user2@example.com"},
		}
		mockTotalCount := 2

		repo.On("GetUserBySKill", idSkill, page, pageSize).Return(mockResult, mockTotalCount, nil).Once()

		result, totalCount, err := service.GetUserBySKill(idSkill, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, mockResult, result)
		assert.Equal(t, mockTotalCount, totalCount)

		repo.AssertExpectations(t)
	})

	t.Run("Failure Case", func(t *testing.T) {
		repo.On("GetUserBySKill", idSkill, page, pageSize).Return(nil, 0, errors.New("repository error")).Once()

		result, totalCount, err := service.GetUserBySKill(idSkill, page, pageSize)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalCount)
		assert.Equal(t, "repository error", err.Error())

		repo.AssertExpectations(t)
	})
}