package event

import (
	"gamma/app/api/core"
	"gamma/app/services/event"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetBootstrapData(c echo.Context) error {
	uuid := c.Get("_id").(string)
	user, err := event.EventRepo().GetUserByUUID(c.Request().Context(), uuid)

	if err != nil {
		panic("Bootstrap < Tailwind")
	}

	join, err := event.EventRepo().GetOrganizationsByUserID(c.Request().Context(), user.UserID)

	if err != nil {
		panic("Uh oh spaghetti-o's")
	}

	return c.JSON(http.StatusOK, core.ApiConverter(map[string]interface{}{
		"bootstrap": join,
	}, 0))
}

func GetAttendingEvents(c echo.Context) error {
	uuid := c.Get("_id").(string)
	user, err := event.EventRepo().GetUserByUUID(c.Request().Context(), uuid)

	if err != nil {
		panic("Yeah idk how this is messing up lol")
	}

	attendingEvents, err := event.EventRepo().GetUserAttendingEvents(c.Request().Context(), user.UserID)

	if err != nil {
		panic("Dude I'm really struggling to understand how this is an error")
	}

	return c.JSON(http.StatusOK, core.ApiConverter(map[string]interface{}{
		"events": attendingEvents,
	}, 0))
}

func GetEventApplications(c echo.Context) error {
	uuid := c.Get("_id").(string)
	user, err := event.EventRepo().GetUserByUUID(c.Request().Context(), uuid)

	if err != nil {
		panic("Y'know writing these panics is really annoying")
	}

	eventApplications, err := event.EventRepo().GetUserEventApplications(c.Request().Context(), user.UserID)

	if err != nil {
		panic("Damn! Even this is panicking???")
	}

	return c.JSON(http.StatusOK, core.ApiConverter(map[string]interface{}{
		"applications": eventApplications,
	}, 0))
}
