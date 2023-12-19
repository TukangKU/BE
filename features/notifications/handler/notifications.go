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

func (nh *notifHandler) GetNotifs() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := jwt.ExtractToken(c.Get("user").(*golangjwt.Token))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "please login first",
			})
		}
		result, err := nh.service.GetNotifs(userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "please login first",
			})
		}
		if result == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "no notification for this user",
				"data":    nil})
		}

		var response = new([]NotifResponse)
		for _, element := range result {
			var notifResp = new(NotifResponse)
			notifResp.ID = element.ID
			notifResp.Message = element.Message
			notifResp.CreatedAt = element.CreatedAt
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success getting notifs",
			"data":    response,
		})
	}
}
