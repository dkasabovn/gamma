package user

import (
	"gamma/app/api/models/auth"

	"github.com/labstack/echo/v4"
)

func loginController(c echo.Context) error {
	var rawLogin auth.UserLogin
	if err := c.Bind(rawLogin); err != nil {
		return err
	}

	return nil
}
