package user_api

import (
	"fmt"
	"gamma/app/api/core"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (a *UserAPI) getUserController(c echo.Context) error {
	fmt.Print(c.Get(echo.HeaderAuthorization))
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*ecJwt.GammaClaims)
	uuid := claims.Uuid

	user, err := a.srvc.GetUser(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"user" : user,
	}))


}
