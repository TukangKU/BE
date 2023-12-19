package main

import (
	"tukangku/config"
	jh "tukangku/features/jobs/handler"
	jr "tukangku/features/jobs/repository"
	js "tukangku/features/jobs/services"
	sr "tukangku/features/skill/repository"
	uh "tukangku/features/users/handler"
	ur "tukangku/features/users/respository"
	us "tukangku/features/users/services"
	ek "tukangku/helper/enkrip"
	"tukangku/routes"
	cld "tukangku/utils/cloudinary"
	"tukangku/utils/database"

	"github.com/labstack/echo/v4"

	nh "tukangku/features/notifications/handler"
	nr "tukangku/features/notifications/repository"
	ns "tukangku/features/notifications/services"
	sh "tukangku/features/skill/handler"
	ss "tukangku/features/skill/services"
)

func main() {

	e := echo.New()

	cfg := config.InitConfig()

	if cfg == nil {
		e.Logger.Fatal("tidak bisa start karena ENV error")
		return
	}
	cld, ctx, param := cld.InitCloudnr(*cfg)

	db, err := database.InitMySQL(*cfg)
	if err != nil {
		e.Logger.Fatal("tidak bisa start bro", err.Error())
	}
	db.AutoMigrate(jr.JobModel{}, ur.UserModel{}, sr.SkillModel{}, nr.NotifModel{})

	// config users features
	enkrip := ek.New()
	userRepo := ur.New(db)
	userService := us.New(userRepo, enkrip)
	userHandler := uh.New(userService, cld, ctx, param)

	// config skill
	skillRepo := sr.New(db)
	skillService := ss.New(skillRepo)
	skillHandler := sh.New(skillService)
	// config jobs
	jobRepo := jr.New(db)
	jobServices := js.New(jobRepo)
	jobHandler := jh.New(jobServices)

	// config notifs
	notifRepo := nr.New(db)
	notifServices := ns.New(notifRepo)
	notifHandler := nh.New(notifServices)

	routes.InitRute(e, userHandler, skillHandler, jobHandler, notifHandler)
	e.Logger.Fatal(e.Start(":8000"))

}
