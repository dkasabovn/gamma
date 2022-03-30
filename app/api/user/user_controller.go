package user_api

import (
	"gamma/app/api/core"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (a *UserAPI) getUserController(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*ecJwt.GammaClaims)
	uuid := claims.Uuid

	user, err := a.srvc.GetUser(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"user": user,
	}))
}

func GetEventsController(c echo.Context) error {
	userObj, err := core.ExtractUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	events, err := user.GetUserService().GetUserOrganizations(c.Request().Context(), userObj.Id)

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"events": events,
	}))
}

func InsertEventController(c echo.Context) error {
	userObj, err := core.ExtractUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	var event bo.Event
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	eventCreated, err := user.GetUserService().InsertEventByOrganization()
}