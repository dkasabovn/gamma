package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"gamma/app/system/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getSelfController(c echo.Context) error {
	self, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	userEvents, err := user.GetUserService().GetUserEvents(c.Request().Context(), self.ID)
	if err != nil {
		log.Errorf("could not get user events: %v", err)
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	userOrgs, err := user.GetUserService().GetUserOrganizations(c.Request().Context(), self.ID)
	if err != nil {
		log.Errorf("could nto get user organizations: %v", err)
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"user":               dto.ConvertUser(self),
		"user_events":        dto.ConvertUserEvents(userEvents),
		"user_organizations": dto.ConvertOrganizationsWithPermissions(userOrgs),
	}))
}
