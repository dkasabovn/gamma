package user

import (
	"net/http"

	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"gamma/app/system/log"

	"github.com/google/uuid"
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
		return c.JSON(http.StatusBadRequest, core.ApiSuccess(map[string]interface{}{
			"error": "1",
		}))
	}

	orgUser, err := core.ExtractOrguser(c, eventCreateDto.OrganizationID)
	if err != nil {
		// not working for some reason
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	// if err := core.FormImage(c, eventCreateDto.EventImage, "event_image"); err != nil {
	// 	return c.JSON(http.StatusBadRequest, core.ApiSuccess(map[string]interface{}{
	// 		"error": "2",
	// 	}))
	// }

	if err := user.GetUserService().CreateEvent(c.Request().Context(), orgUser, &eventCreateDto); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiSuccess(map[string]interface{}{
			"error": "3",
		}))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}

func checkController(c echo.Context) error {
	var eventCheckDto dto.EventCheck
	if err := c.Bind(&eventCheckDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	userEvents, err := user.GetUserService().GetUserEvents(c.Request().Context(), uuid.Must(uuid.Parse(eventCheckDto.UserID)))
	if err != nil {
		log.Errorf("could not get user events: %v", err)
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	for i := 0; i < len(userEvents); i++ {
		if userEvents[i].EventFk == uuid.Must(uuid.Parse(eventCheckDto.EventID)) {
			return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
		}
	}
	return core.JSONApiError(c, http.StatusUnauthorized)

}
