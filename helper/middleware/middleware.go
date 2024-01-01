package middleware

import (
	"net/http"
	"tukangku/helper/jwt"
	"tukangku/helper/responses"

	golangjwt "github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CheckClient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole, _ := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))

		if userRole != "client" {
			return responses.PrintResponse(c, http.StatusUnauthorized, "anda bukan client", nil)
		}
		return next(c)
	}
}

func CheckWorker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole, _ := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))

		if userRole != "worker" {
			return responses.PrintResponse(c, http.StatusUnauthorized, "anda bukan worker", nil)
		}
		return next(c)
	}
}

func LogMiddleware(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "\nhost=${host}, uri=${uri}, user_agent=${user_agent}, time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))
}
