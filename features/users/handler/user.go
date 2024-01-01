package user

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"strings"
	"tukangku/features/skill"
	"tukangku/features/users"
	"tukangku/helper/jwt"
	"tukangku/helper/responses"

	cld "tukangku/utils/cloudinary"

	golangjwt "github.com/golang-jwt/jwt/v5"

	"github.com/cloudinary/cloudinary-go/v2"

	echo "github.com/labstack/echo/v4"
)

type userController struct {
	srv    users.Service
	cl     *cloudinary.Cloudinary
	ct     context.Context
	folder string
}

func New(s users.Service, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) users.Handler {
	return &userController{
		srv:    s,
		cl:     cld,
		ct:     ctx,
		folder: uploadparam,
	}
}

func (ur *userController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(RegisterRequest)
		if err := c.Bind(input); err != nil {
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"input yang di berikan tidak sesuai",
				nil)
		}
		var inputProcess = new(users.Users)
		inputProcess.UserName = input.UserName
		inputProcess.Email = input.Email
		inputProcess.Password = input.Password
		inputProcess.Role = input.Role

		result, err := ur.srv.Register(*inputProcess)

		if err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "harus di isi") {
				return responses.PrintResponse(
					c, http.StatusBadRequest,
					strings.ReplaceAll(err.Error(), "", ""),
					nil)
			}
			if strings.Contains(err.Error(), "duplicate entry") {
				return responses.PrintResponse(
					c, http.StatusConflict,
					"data sudah ada",
					nil)
			}
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"internal server error",
				nil)
		}

		var responsess = new(RegisterResponse)
		responsess.ID = result.ID
		responsess.UserName = result.UserName
		responsess.Email = result.Email
		responsess.Role = result.Role

		return responses.PrintResponse(
			c, http.StatusCreated,
			"success create data",
			responsess)
	}
}

func (ul *userController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoginRequest)
		if err := c.Bind(input); err != nil {
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"input yang diberikan tidak sesuai",
				nil)
		}

		result, err := ul.srv.Login(input.Email, input.Password)

		if err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "harus di isi") {
				return responses.PrintResponse(
					c, http.StatusBadGateway,
					strings.ReplaceAll(err.Error(), "", ""),
					nil)
			}

			if strings.Contains(err.Error(), "record not found") {
				return responses.PrintResponse(
					c, http.StatusNotFound,
					"record not found",
					nil)
			}

			if strings.Contains(err.Error(), "password salah") {
				return responses.PrintResponse(
					c, http.StatusUnauthorized,
					"password salah",
					nil)
			}

			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"internal server error",
				nil)

		}

		strToken, err := jwt.GenerateJWT(result.ID, result.Role)
		if err != nil {
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"terjadi permasalahan ketika mengenkripsi data",
				nil)
		}

		var response = new(LoginResponse)
		response.ID = result.ID
		response.UserName = result.UserName
		response.Email = result.Email
		response.Role = result.Role
		response.Token = strToken

		
		return responses.PrintResponse(
			c, http.StatusOK,
			"success create data",
			response)
	}
}

func (us *userController) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		var input = new(UserUpdate)
		if err := c.Bind(input); err != nil {
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"invalid input",
				nil)
		}

		formHeader, _ := c.FormFile("foto")

		var link string

		if formHeader != nil {

			formFile, err := formHeader.Open()
			if err != nil {
				return responses.PrintResponse(
					c, http.StatusInternalServerError,
					"formfile error",
					nil)
			}

			link, err = cld.UploadImage(us.cl, us.ct, formFile, us.folder)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					return responses.PrintResponse(
						c, http.StatusBadRequest,
						"harap pilih gambar",
						nil)
				} else {
					return responses.PrintResponse(
						c, http.StatusInternalServerError,
						"kesalahan pada server",
						nil)
				}
			}

			input.Foto = link

		}

		updatedCient := users.Users{
			UserName: input.UserName,
			Nama:     input.Nama,
			Email:    input.Email,
			NoHp:     input.NoHp,
			Alamat:   input.Alamat,
			Foto:     input.Foto,
			ID:       userID,
		}
		for _, v := range input.Skills {
			updatedCient.Skill = append(updatedCient.Skill, skill.Skills{ID: v})
		}

		result, err := us.srv.UpdateUser(userID, updatedCient)
		if err != nil {
			c.Logger().Error("ERROR UpdateUser, explain:", err.Error())
			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"failed to update user",
				nil)
		}

		result.Foto = link

		var response = new(UserResponseUpdate)
		response.ID = result.ID
		response.UserName = result.UserName
		response.Nama = result.Nama
		response.Email = result.Email
		response.NoHp = result.NoHp
		response.Alamat = result.Alamat
		response.Foto = result.Foto
		response.ID = userID
		response.Skill = func() []UserSkill {
			var skill []UserSkill
			for _, s := range result.Skill {
				skill = append(skill, UserSkill{
					SkillID:   s.ID,
					NamaSKill: s.NamaSkill,
				})
			}
			return skill
		}()

		return responses.PrintResponse(
			c, http.StatusOK,
			"posting update successfully",
			response)
	}
}

func (gu *userController) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"invalid ID",
				nil)
		}

		results, err := gu.srv.GetUserByID(uint(userID))
		if err != nil {
			c.Logger().Error("ERROR GetByID, explain:", err.Error())

			if strings.Contains(err.Error(), "not found") {
				return responses.PrintResponse(
					c, http.StatusNotFound,
					"posting not found",
					nil)
			}

			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"error retrieving posting by id",
				nil)
		}
		response := UserResponseUpdate{
			ID:       results.ID,
			UserName: results.UserName,
			Nama:     results.Nama,
			Email:    results.Email,
			NoHp:     results.NoHp,
			Alamat:   results.Alamat,
			Role:     results.Role,
			Foto:     results.Foto,
			Skill: func() []UserSkill {
				var skill []UserSkill
				for _, s := range results.Skill {
					skill = append(skill, UserSkill{
						SkillID:   s.ID,
						NamaSKill: s.NamaSkill,
					})
				}
				return skill
			}(),
			Job: func() []UserJob {
				var skill []UserJob
				for _, s := range results.Job {
					skill = append(skill, UserJob{
						JobID:    s.ID,
						Price:    s.Price,
						Category: s.Category,
					})
				}
				return skill
			}(),
			JobCount: results.JobCount,
		}

		return responses.PrintResponse(
			c, http.StatusOK,
			"success get data by id",
			response)
	}
}

func (gu *userController) GetUserBySKill() echo.HandlerFunc {
	return func(c echo.Context) error {
		skillID, err := strconv.Atoi(c.QueryParam("skill"))
		if err != nil {
			return responses.PrintResponse(
				c, http.StatusBadRequest,
				"invalid skill id",
				nil)
		}

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.QueryParam("pagesize"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}

		users, totalCount, err := gu.srv.GetUserBySKill(uint(skillID), page, pageSize)
		if err != nil {
			c.Logger().Error("ERROR GetUserBySkill, explain:", err.Error())

			return responses.PrintResponse(
				c, http.StatusInternalServerError,
				"error retrieving users by skill",
				nil)
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

		var responsess []GetUserResponse
		for _, v := range users {
			responsess = append(responsess, GetUserResponse{
				ID:       v.ID,
				Nama:     v.Nama,
				UserName: v.UserName,
				Alamat:   v.Alamat,
				Email:    v.Email,
				Foto:     v.Foto,
				Skill: func() []UserSkill {
					var skill []UserSkill
					for _, s := range v.Skill {
						skill = append(skill, UserSkill{
							SkillID:   s.ID,
							NamaSKill: s.NamaSkill,
						})
					}
					return skill
				}(),
				JobCount: v.JobCount,
			})

		}

		return responses.PrintResponse(
			c, http.StatusOK,
			"success get data by id",
			map[string]interface{}{
				"data": responsess,
				"pagination": map[string]interface{}{
					"page":       page,
					"pagesize":   pageSize,
					"totalpages": totalPages,
				},
			})
	}
}
