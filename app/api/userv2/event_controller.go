package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"gamma/app/system/log"
	"net/http"

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

func validateOtherController(c echo.Context) error {
	var eventValidateDto dto.EventValidate
	if err := c.Bind(&eventValidateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	_, err := core.ExtractOrguser(c, eventValidateDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	// TODO: Check if user has user_event matching the requested validation, ensure orguser has the right privileges, send back users information

	return nil
}

func getEventController(c echo.Context) error {
	stringID := c.Param("event_uuid")
	eID, err := uuid.Parse(stringID)
	if err != nil {
		log.Errorf("incorrect uuid: %v", err)
	}
	_, err = core.ExtractUser(c)
	if err != nil {
		log.Errorf("extract user: %v", err)
		return core.JSONApiError(c, http.StatusUnauthorized)
	}
	event, err := user.GetUserService().GetEvent(c.Request().Context(), eID)

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"event": dto.ConvertEvent(event),
	}))
}

func updateEventController(c echo.Context) error {
	stringID := c.Param("event_uuid")
	eID, err := uuid.Parse(stringID)
	if err != nil {
		log.Errorf("incorrect uuid: %v", err)
	}
	_, err = core.ExtractUser(c)
	if err != nil {
		log.Errorf("extract user: %v", err)
		return core.JSONApiError(c, http.StatusUnauthorized)
	}
	var eventCreateDto dto.EventUpsert
	if err := c.Bind(&eventCreateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	orgUser, err := core.ExtractOrguser(c, eventCreateDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	// if err := core.FormImage(c, eventCreateDto.EventImage, "event_image"); err != nil {
	// 	return core.JSONApiError(c, http.StatusBadRequest)
	// }

	err = user.GetUserService().UpdateEvent(c.Request().Context(), orgUser, &eventCreateDto, eID)
	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"name":      "scooze",
		"eventname": eventCreateDto.EventName,
	}))
}
