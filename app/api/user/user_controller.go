package user_api

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
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

func (a *UserAPI) getUserOrganizationsController(c echo.Context) error {
	org_user, err := core.ExtractOrguser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	orgs, err := a.srvc.GetUserOrganizations(c.Request().Context(), org_user.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, core.ApiError(http.StatusNotFound))
	}

	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"organizations": dto.ConvertOrganizations(orgs),
	}))
}

func (a *UserAPI) getEventsController(c echo.Context) error {
	// _, err := core.ExtractUser(c)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	// }
	events, err := a.srvc.GetEvents(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}
	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"events": dto.ConvertEvents(events),
	}))
}
