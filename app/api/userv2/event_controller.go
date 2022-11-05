package user

import (
	"net/http"

	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"gamma/app/system/log"

	"github.com/labstack/echo/v4"
)

func getEventsController(c echo.Context) error {
	_, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("extract user: %v", err)
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	var eventSearchDto dto.EventSearch
	if err := c.Bind(&eventSearchDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	events, err := user.GetUserService().GetEvents(c.Request().Context(), &eventSearchDto)
	if err != nil {
		log.Errorf("getting events: %v", err)
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
		//status unauthorized
		return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
			"reason": "issue 1",
		}))
	}

	if err := core.FormImage(c, eventCreateDto.EventImage, "event_image"); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	if err := user.GetUserService().CreateEvent(c.Request().Context(), orgUser, &eventCreateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}

func checkController(c echo.Context) error {
	var eventCheckDto dto.EventCheck
	if err := c.Bind(&eventCheckDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	if err := user.GetUserService().CheckUser(c.Request().Context(), eventCheckDto.UserID, eventCheckDto.EventID); err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}
