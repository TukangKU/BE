package routes

import (
	"net/http"
	"tukangku/features/jobs"

	"tukangku/features/skill"
	"tukangku/features/transaction"
	"tukangku/features/users"
	"tukangku/helper/jwt"
	"tukangku/helper/responses"

	golangjwt "github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRute(e *echo.Echo, uh users.Handler, us skill.Handler, jh jobs.Handler, th transaction.Handler) {

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	routeUser(e, uh)
	routeSkill(e, us)
	routeJobs(e, jh)

	routeTransaction(e, th)

}

func routeUser(e *echo.Echo, uh users.Handler) {

	e.POST("/register", uh.Register())
	e.POST("/login", uh.Login())
	e.PUT("/client/:id", uh.UpdateUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")), checkClient)
	e.PUT("/worker/:id", uh.UpdateUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")), checkWorker)
	e.GET("client/:id", uh.GetUserByID(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("worker/:id", uh.GetUserByID(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/users/skill", uh.GetUserBySKill(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/takeworker", uh.TakeWorker(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func routeSkill(e *echo.Echo, us skill.Handler) {
	e.POST("/addskill", us.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/skills", us.Show(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func routeJobs(e *echo.Echo, jh jobs.Handler) {
	e.POST("/jobs", jh.Create(), echojwt.JWT([]byte("$!1gnK3yyy!!!")), checkClient)
	e.GET("/jobs", jh.GetJobs(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/jobs/:id", jh.GetJob(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PUT("/jobs/:id", jh.UpdateJob(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

}

func routeTransaction(e *echo.Echo, th transaction.Handler) {
	e.POST("/transaction", th.AddTransaction(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/transaction/:id", th.CheckTransaction(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.POST("/callback", th.CallBack())
}

func checkClient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Your middleware logic here
		userRole, _ := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))
		// Access request path with c.Path()

		// Example: Only execute middleware for paths starting with "/api"
		if userRole != "client" {
			// c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusUnauthorized
			var message = "anda bukan client"

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		return next(c)
	}
}
func checkWorker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Your middleware logic here
		userRole, _ := jwt.ExtractTokenRole(c.Get("user").(*golangjwt.Token))
		// Access request path with c.Path()

		// Example: Only execute middleware for paths starting with "/api"
		if userRole != "worker" {
			// c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusUnauthorized
			var message = "anda bukan worker"

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		return next(c)
	}
}
