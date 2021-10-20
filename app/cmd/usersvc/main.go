package main

import (
	"gamma/app/api/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	// adds temp get and post routes
	user.JwtRoutes(e)

	// temp no auth rout
	e.GET("/", noAuth)

	e.Logger.Fatal(e.Start(":8000"))

}

func noAuth(ctx echo.Context) error {
	return ctx.JSON(http.StatusAccepted, "No auth needed")
}
