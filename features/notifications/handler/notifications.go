package notifications

import (
	"net/http"
	"tukangku/features/notifications"
	"tukangku/helper/jwt"

	golangjwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type notifHandler struct {
	service notifications.Service
}

func New(s notifications.Service) notifications.Handler {
	return &notifHandler{
		service: s,
	}
}

func (sh *notifHandler) GetNotifs() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "please login first",
			})
		}

	}
}
