package main

import (
	"gamma/app/api/event"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	event.SetUpInfoGroup(e)
	e.Logger.Fatal(e.Start(":8080"))
}
