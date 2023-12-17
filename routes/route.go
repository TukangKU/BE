package routes

import (
	"tukangku/features/users"
	"tukangku/features/skill"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRute(e *echo.Echo, uh users.Handler, us skill.Handler) {

	// e.Use(mdd.CORS())
	// e.Use(mdd.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())


	routeUser(e, uh)
	routeSkill(e, us)


}

func routeUser(e *echo.Echo, uh users.Handler) {
	e.POST("/register", uh.Register())
	e.POST("/login", uh.Login())
	e.PUT("/client/:id", uh.UpdateUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PUT("/worker/id", uh.UpdateUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func routeSkill(e *echo.Echo, us skill.Handler) {
	e.POST("/addskill", us.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/skills", us.Show(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}
