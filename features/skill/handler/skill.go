package skill

import (
	"net/http"
	"tukangku/features/skill"
	"tukangku/helper/responses"

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

		if err := c.Bind(&inputSkill); err != nil {
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"invalid input",
				nil)
		}

		var inputProcess = new(skill.Skills)
		inputProcess.NamaSkill = inputSkill.NamaSkill

		result, err := sh.s.AddSkill(*inputProcess)
		if err != nil {
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"failed to add skill",
				nil)
		}

		return responses.PrintResponse(
			c, http.StatusCreated,
			"success adding skill",
			result)
	}
}

func (sh *SkillHandler) Show() echo.HandlerFunc {
	return func(c echo.Context) error {
		skills, err := sh.s.ShowSkill()
		if err != nil {
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"failed to retrieve skills",
				nil)
		}

		return responses.PrintResponse(
			c, http.StatusOK,
			"success retrieving skill",
			skills)
	}
}
