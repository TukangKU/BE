package main

import (
	"tukangku/config"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	cfg := config.InitConfig()

	if cfg == nil {
		e.Logger.Fatal("tidak bisa start karena ENV error")
		return
	}

	e.Logger.Fatal(e.Start(":8000"))

}