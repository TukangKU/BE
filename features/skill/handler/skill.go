package skill

import (
	"net/http"
	"tukangku/features/skill"

	// golangjwt "github.com/golang-jwt/jwt/v5"

	echo "github.com/labstack/echo/v4"
)

type SkillHandler struct {
	s skill.Service
}

func New(s skill.Service) skill.Handler {
	return &SkillHandler{
		s: s,
	}
}

func (sh *SkillHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var inputSkill = new(SkillRequest)

		// Bind the request body to the skill model
		if err := c.Bind(&inputSkill); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid input",
			})
		}

		var inputProcess = new(skill.Skills)
		inputProcess.NamaSkill = inputSkill.NamaSkill

		// Call the service method to add the skill
		result, err := sh.s.AddSkill(*inputProcess)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "failed to add skill",
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "success adding skill",
			"data":    result,
		})
	}
}

func (sh *SkillHandler) Show() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Call the service method to retrieve all skills
		skills, err := sh.s.ShowSkill()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "failed to retrieve skills",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success retrieving skills",
			"data":    skills,
		})
	}
}
