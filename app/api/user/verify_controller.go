package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getUser(c echo.Context) error {
	return c.JSON(http.StatusContinue, "no users yet!")
}

func createUser(c echo.Context) error {
	return c.JSON(http.StatusContinue, "no users yet!")
}