package user

import (
	"context"
	"fmt"
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
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang di berikan tidak sesuai",
			})
		}
		var inputProcess = new(users.Users)
		inputProcess.UserName = input.UserName
		inputProcess.Email = input.Email
		inputProcess.Password = input.Password
		inputProcess.Role = input.Role

		if inputProcess.UserName == "" {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "masukkan username",
			})
		}

		if inputProcess.Email == "" {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "masukkan email",
			})
		}

		if inputProcess.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "masukkan password",
			})
		}

		if inputProcess.Role == "" {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "masukkan role",
			})
		}

		result, err := ur.srv.Register(*inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "terdaftar") {
				statusCode = http.StatusBadRequest
				message = "data yang diinputkan sudah terdaftar ada sistem"
			}

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		var responsess = new(RegisterResponse)
		responsess.ID = result.ID
		responsess.UserName = result.UserName
		responsess.Email = result.Email
		responsess.Role = result.Role

		return responses.PrintResponse(c, http.StatusCreated, "success create data", responsess)
	}
}

func (ul *userController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoginRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		result, err := ul.srv.Login(input.Email, input.Password)

		if err != nil {
			c.Logger().Error("ERROR Login, explain:", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "data yang diinputkan tidak ditemukan",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		strToken, err := jwt.GenerateJWT(result.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika mengenkripsi data",
			})
		}

		var responses = new(LoginResponse)
		responses.ID = result.ID
		responses.UserName = result.UserName
		responses.Email = result.Email
		responses.Role = result.Role
		responses.Token = strToken

		return c.JSON(http.StatusOK, map[string]any{
			"message": "success create data",
			"data":    responses,
		})
	}
}

func (us *userController) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		var input = new(UserUpdate)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"massage": "invalid input",
			})
		}

		formHeader, _ := c.FormFile("foto")

		var link string

		if formHeader != nil {

			formFile, err := formHeader.Open()
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError, map[string]any{
						"message": "formfile error",
					})
			}

			link, err = cld.UploadImage(us.cl, us.ct, formFile, us.folder)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					return c.JSON(http.StatusBadRequest, map[string]any{
						"message": "harap pilih gambar",
						"data":    nil,
					})
				} else {
					return c.JSON(http.StatusInternalServerError, map[string]any{
						"message": "kesalahan pada server",
						"data":    nil,
					})
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
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "failed to update user",
			})
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

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "posting updated successfully",
			"data":    response,
		})
	}
}

func (gu *userController) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))

		results, err := gu.srv.GetUserByID(uint(userID))
		if err != nil {
			c.Logger().Error("ERROR GetByID, explain:", err.Error())

			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"message": "Posting not found",
				})
			}

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error retrieving Posting by ID",
			})
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
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get data by ID",
			"data":    response,
		})
	}
}

func (gu *userController) GetUserBySKill() echo.HandlerFunc {
	return func(c echo.Context) error {
		skillID, err := strconv.Atoi(c.QueryParam("skill"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid skill id",
			})
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

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error retrieving users by skill",
			})
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

		var responsess []GetUserResponse
		for _, v := range users {
			responsess = append(responsess, GetUserResponse{
				ID:       v.ID,
				Nama:     v.Nama,
				UserName: v.UserName,
				Alamat:   v.Alamat,
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
			})

		}

		fmt.Println("response user =", responsess)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get data by ID",
			"data":    responsess,
			"pagination": map[string]interface{}{
				"skill":      skillID,
				"page":       page,
				"pagesize":   pageSize,
				"totalPages": totalPages},
		})
	}
}
