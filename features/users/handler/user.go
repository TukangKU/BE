package user

import (
	"context"
	"net/http"
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

// Register implements users.Handler.
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

// Login implements users.Handler.
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

// UpdateUser implements users.Handler.
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

		var response = &UserUpdate{
			UserName: result.UserName,
			Nama:     result.Nama,
			Email:    result.Email,
			NoHp:     result.NoHp,
			Alamat:   result.Alamat,
			Foto:     result.Foto,
			Role:     result.Role,
			ID:       userID,
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "posting updated successfully",
			"data":    response,
		})
	}
}
