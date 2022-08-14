package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getEventsController(c echo.Context) error {
	_, err := core.ExtractUser(c)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	var eventSearchDto dto.EventSearch
	if err := c.Bind(&eventSearchDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	events, err := user.GetUserService().GetEvents(c.Request().Context(), &eventSearchDto)
	if err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"events": dto.ConvertEventsWithOrganizations(events),
	}))
}

func createEventController(c echo.Context) error {
	var eventCreateDto dto.EventUpsert
	if err := c.Bind(&eventCreateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	orgUser, err := core.ExtractOrguser(c, eventCreateDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	if err := core.FormImage(c, eventCreateDto.EventImage, "event_image"); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	if err := user.GetUserService().CreateEvent(c.Request().Context(), orgUser, &eventCreateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}
