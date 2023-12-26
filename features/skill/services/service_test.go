package services_test

import (
	"errors"
	"testing"
	"tukangku/features/skill"
	"tukangku/features/skill/mocks"
	"tukangku/features/skill/services"

	"github.com/stretchr/testify/assert"
)

func TestShowSkill(t *testing.T) {
	repo := mocks.NewRepository(t)
	skillService := services.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockSkills := []skill.Skills{
			{ID: 1, NamaSkill: "Service Ac"},
			{ID: 2, NamaSkill: "Plumber"},
		}

		repo.On("ShowSkill").Return(mockSkills, nil).Once()

		result, err := skillService.ShowSkill()

		assert.Nil(t, err)
		assert.Equal(t, mockSkills, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		repo.On("ShowSkill").Return(nil, errors.New("failed to retrieve skills")).Once()

		_, err := skillService.ShowSkill()

		assert.Error(t, err)
		assert.Equal(t, "failed to retrieve skills", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestAddSkill(t *testing.T) {
	repo := mocks.NewRepository(t)
	skillService := services.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockNewSkill := skill.Skills{NamaSkill: "Plumber"}

		repo.On("AddSkill", mockNewSkill).Return(skill.Skills{ID: 1, NamaSkill: "Plumber"}, nil).Once()

		result, err := skillService.AddSkill(mockNewSkill)

		assert.Nil(t, err)
		assert.Equal(t, skill.Skills{ID: 1, NamaSkill: "Plumber"}, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockNewSkill := skill.Skills{NamaSkill: "Plumber"}

		repo.On("AddSkill", mockNewSkill).Return(skill.Skills{}, errors.New("failed to add skill")).Once()

		_, err := skillService.AddSkill(mockNewSkill)

		assert.Error(t, err)
		assert.Equal(t, "failed to add skill", err.Error())

		repo.AssertExpectations(t)
	})
}