package main

import (
	"gamma/app/api/user"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	user.AddOpenRoutes(e)

	e.Logger.Fatal(e.Start(":6969"))
}
