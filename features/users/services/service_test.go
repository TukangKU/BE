package services_test

import (
	"errors"
	"testing"
	"tukangku/features/users"
	"tukangku/features/users/mocks"
	"tukangku/features/users/services"
	eMock "tukangku/helper/enkrip/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := eMock.NewHashInterface(t)

	s := services.New(repo, enkrip)

	var inputData = users.Users{
		UserName: "bedul",
		Email:    "bedul@gmail.com",
		Role:     "client",
		Password: "bedul",
	}

	var repoData = users.Users{
		UserName: "bedul",
		Email:    "bedul@gmail.com",
		Role:     "client",
		Password: "some string",
	}

	var successReturnData = users.Users{
		ID:       uint(1),
		UserName: "bedul",
		Email:    "bedul@gmail.com",
		Role:     "client",
	}

	// var fasleData = users.Users{}

	t.Run("Success Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("some string", nil).Once()

		repo.On("Register", repoData).Return(successReturnData, nil).Once()
		res, err := s.Register(inputData)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, "", res.Password)
	})

	// t.Run("Failed Case", func(t *testing.T) {
	// 	s := services.New(&FalseMockRepository{}, &MockEncryp{})
	// 	res, err := s.Register(fasleData)
	// 	assert.Error(t, err)
	// 	assert.Equal(t, uint(0), res.ID)
	// 	assert.Equal(t, "", res.Email)

	// })

	t.Run("Hashing Error Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("", errors.New("hashing error")).Once()

		newUser := users.Users{
			Email:    "bedul@gmail.com",
			Password: "bedul",
		}
		result, err := s.Register(newUser)
		enkrip.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "terdapat masalah saat memproses data", err.Error())

	})

	t.Run("Duplicate Error Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("bedul", nil).Once()

		newUser := users.Users{
			Email:    "bedul@gmail.com",
			Password: "bedul",
		}
		repo.On("Register", newUser).Return(users.Users{}, errors.New("duplicate")).Once()
		result, err := s.Register(newUser)
		enkrip.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "data telah terdaftar pada sistem", err.Error())

	})

	t.Run("Repository Error Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("bedul", nil).Once()

		newUser := users.Users{
			Email:    "bedul@gmail.com",
			Password: "bedul",
		}
		repo.On("Register", newUser).Return(users.Users{}, errors.New("bedul")).Once()
		result, err := s.Register(newUser)
		enkrip.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())

	})

	t.Run("Failure Case - Incorrect Input Data", func(t *testing.T) {
		inputData := users.Users{}

		result, err := s.Register(inputData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "incorrect inpu data", err.Error())
	})

}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	hashMock := eMock.NewHashInterface(t)

	service := services.New(repo, hashMock)

	// Success Case
	t.Run("Success Case", func(t *testing.T) {
		email := "test@example.com"
		password := "testpassword"

		// Mocking the Login call
		repo.On("Login", email).Return(users.Users{Password: "hashedpassword"}, nil).Once()

		hashMock.On("Compare", "hashedpassword", password).Return(nil).Once()

		result, err := service.Login(email, password)

		hashMock.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, "hashedpassword", result.Password)
	})

	t.Run("Failure Case - Incorrect Input Data", func(t *testing.T) {
		result, err := service.Login("", "testpassword")

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "incorrect input data", err.Error())
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		email := "nonexistent@example.com"

		repo.On("Login", email).Return(users.Users{}, errors.New("not found")).Once()

		result, err := service.Login(email, "testpassword")

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "data tidak ditemukan", err.Error())
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		email := "nonexistent@example.com"

		repo.On("Login", email).Return(users.Users{}, errors.New("bedul")).Once()

		result, err := service.Login(email, "testpassword")

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
	})

	t.Run("Failure Case - Incorrect Password", func(t *testing.T) {
		email := "test@example.com"
		password := "incorrectpassword"

		repo.On("Login", email).Return(users.Users{Password: "hashedpassword"}, nil).Once()

		hashMock.On("Compare", "hashedpassword", password).Return(errors.New("password salah")).Once()

		result, err := service.Login(email, password)

		hashMock.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "password salah", err.Error())
	})
}

func TestGetUserByID(t *testing.T) {
	var userID = uint(1)
	repo := mocks.NewRepository(t)
	enkrip := eMock.NewHashInterface(t)
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

func TestUpdateUser(t *testing.T) {
	repo := mocks.NewRepository(t)
	service := services.New(repo, nil)

	idUser := uint(1)
	updateData := users.Users{
		UserName: "updatedUser",
		Email:    "updateduser@example.com",
	}

	t.Run("Success Case", func(t *testing.T) {
		mockResult := users.Users{
			ID:       idUser,
			UserName: "updatedUser",
			Email:    "updateduser@example.com",
		}

		repo.On("UpdateUser", idUser, updateData).Return(mockResult, nil).Once()

		result, err := service.UpdateUser(idUser, updateData)

		assert.NoError(t, err)
		assert.Equal(t, mockResult, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		repo.On("UpdateUser", idUser, updateData).Return(users.Users{}, errors.New("not found")).Once()

		result, err := service.UpdateUser(idUser, updateData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "failed to update user", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - Other Error", func(t *testing.T) {
		repo.On("UpdateUser", idUser, updateData).Return(users.Users{}, errors.New("repository error")).Once()

		result, err := service.UpdateUser(idUser, updateData)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "failed to update user", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestTakeWorker(t *testing.T) {
	repo := mocks.NewRepository(t)
	service := services.New(repo, nil)

	idUser := uint(1)

	t.Run("Success Case", func(t *testing.T) {
		mockResult := users.Users{
			ID:       idUser,
			UserName: "takenUser",
			Email:    "takenuser@example.com",
		}

		repo.On("TakeWorker", idUser).Return(mockResult, nil).Once()

		result, err := service.TakeWorker(idUser)

		assert.NoError(t, err)
		assert.Equal(t, mockResult, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		repo.On("TakeWorker", idUser).Return(users.Users{}, errors.New("not found")).Once()

		result, err := service.TakeWorker(idUser)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "user not found", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - Other Error", func(t *testing.T) {
		repo.On("TakeWorker", idUser).Return(users.Users{}, errors.New("repository error")).Once()

		result, err := service.TakeWorker(idUser)

		assert.Error(t, err)
		assert.Equal(t, users.Users{}, result)
		assert.Equal(t, "error retrieving User by ID", err.Error())

		repo.AssertExpectations(t)
	})
}

type MockRepository struct{}

func (mr *MockRepository) Register(newUser users.Users) (users.Users, error) {
	return users.Users{}, nil
}

func (mr *MockRepository) Login(email string, password string) (users.Users, error) {
	return users.Users{}, nil
}

func (mr *MockRepository) UpdateUser(idUser uint, updateWorker users.Users) (users.Users, error) {
	return users.Users{}, nil
}

func (mr *MockRepository) GetUserByID(idSkill uint) (users.Users, error) {
	return users.Users{}, nil
}

func (mr *MockRepository) GetUserBySKill(idSkill uint, page, pageSize int) ([]users.Users, int, error) {
	return []users.Users{}, 0, nil
}

func (mr *MockRepository) TakeWorker(idUser uint) (users.Users, error) {
	return users.Users{}, nil
}

type FalseMockRepository struct{}

func (mr *FalseMockRepository) Register(newUser users.Users) (users.Users, error) {
	return users.Users{}, errors.New("something wrong")
}

func (mr *FalseMockRepository) Login(email string) (users.Users, error) {
	return users.Users{}, nil
}

func (mr *FalseMockRepository) UpdateUser(idUser uint, updateWorker users.Users) (users.Users, error) {
	return users.Users{}, nil
}

func (mr *FalseMockRepository) GetUserByID(idSkill uint) (users.Users, error) {
	return users.Users{}, nil
}

func (mr *FalseMockRepository) GetUserBySKill(idSkill uint, page, pageSize int) ([]users.Users, int, error) {
	return []users.Users{}, 0, nil
}

func (mr *FalseMockRepository) TakeWorker(idUser uint) (users.Users, error) {
	return users.Users{}, nil
}

type MockEncryp struct{}

func (me *MockEncryp) Compare(hashed string, input string) error {
	return nil
}

func (me *MockEncryp) HashPassword(input string) (string, error) {
	return "bedul", nil
}

// func Test_userService_GetUserByID(t *testing.T) {
// 	type args struct {
// 		idUser uint
// 	}
// 	tests := []struct {
// 		name    string
// 		gu      *userService
// 		args    args
// 		want    users.Users
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := tt.gu.GetUserByID(tt.args.idUser)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("userService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("userService.GetUserByID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
