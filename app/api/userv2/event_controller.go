package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"net/http"

	"github.com/google/uuid"
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

func getEventController(c echo.Context) error {
	_, err := core.ExtractUser(c)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	var eventGetDto dto.EventGet
	if err := c.Bind(&eventGetDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	eventID, err := uuid.Parse(eventGetDto.EventID)
	if err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	event, err := user.GetUserService().GetEvent(c.Request().Context(), eventID)

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"Event": event,
	}))
}

func updateEventController(c echo.Context) error {
	_, err := core.ExtractUser(c)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	var eventUpdateDto dto.EventUpdate
	if err := c.Bind(&eventUpdateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	orgUser, err := core.ExtractOrguser(c, eventUpdateDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}
	eventID, err := uuid.Parse(eventUpdateDto.EventID)
	if err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	event, err := user.GetUserService().GetEvent(c.Request().Context(), eventID)
	if eventUpdateDto.EventName == "" {
		eventUpdateDto.EventName = event.EventName
	}
	if eventUpdateDto.EventDate.IsZero() {
		eventUpdateDto.EventDate = event.EventDate
	}
	if eventUpdateDto.EventDescription == "" {
		eventUpdateDto.EventDescription = event.EventDescription
	}
	if eventUpdateDto.EventLocation == "" {
		eventUpdateDto.EventName = event.EventLocation
	}
	if len(eventUpdateDto.EventImage) == 0 {
		eventUpdateDto.EventImage = []byte(event.EventImageUrl)
	} else {
		if err := core.FormImage(c, eventUpdateDto.EventImage, "event_image"); err != nil {
			return core.JSONApiError(c, http.StatusBadRequest)
		}
	}

	if err := user.GetUserService().UpdateEvent(c.Request().Context(), orgUser, &eventUpdateDto, eventID); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}

func deleteEventController(c echo.Context) error {
	_, err := core.ExtractUser(c)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	var eventDeleteDto dto.EventDelete
	if err := c.Bind(&eventDeleteDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	orgUser, err := core.ExtractOrguser(c, eventDeleteDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	eventID, err := uuid.Parse(eventDeleteDto.EventID)
	if err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	if err := user.GetUserService().DeleteEvent(c.Request().Context(), orgUser, eventID); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))

}
