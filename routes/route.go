package routes

import (
	"tukangku/features/jobs"
	"tukangku/features/notifications"
	"tukangku/features/skill"
	"tukangku/features/users"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRute(e *echo.Echo, uh users.Handler, us skill.Handler, jh jobs.Handler, nh notifications.Handler) {

	// e.Use(mdd.CORS())
	// e.Use(mdd.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	routeUser(e, uh)
	routeSkill(e, us)
	routeJobs(e, jh)
	routeNotifs(e, nh)

}

func routeUser(e *echo.Echo, uh users.Handler) {
	e.POST("/register", uh.Register())
	e.POST("/login", uh.Login())
	e.PUT("/client/:id", uh.UpdateUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PUT("/worker/:id", uh.UpdateUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func routeSkill(e *echo.Echo, us skill.Handler) {
	e.POST("/addskill", us.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/skills", us.Show(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func routeJobs(e *echo.Echo, jh jobs.Handler) {
	e.POST("/jobs", jh.Create(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func routeNotifs(e *echo.Echo, nh notifications.Handler) {
	e.GET("/notifications", nh.GetNotifs(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}
