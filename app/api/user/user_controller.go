package user_api

import (
	"gamma/app/api/core"
	"gamma/app/services/user"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetUserController(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*ecJwt.GammaClaims)
	uuid := claims.Uuid

	user, err := user.GetUserService().GetUser(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"user" : user,
	}))


}
