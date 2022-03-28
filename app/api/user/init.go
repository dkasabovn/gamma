package user

import (
	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()

	AddOpenRoutes(e)
	e.Logger.Fatal(e.Start(":6969"))
}
