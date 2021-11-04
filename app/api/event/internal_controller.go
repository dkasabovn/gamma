package event

import (
	"gamma/app/api/core"
	"gamma/app/datastore/events"
	"gamma/app/datastore/events/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func createUser(c echo.Context) error {
	user := new(models.User)
	err := c.Bind(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(505))
	}
	err = user.Insert(c.Request().Context(), events.EventDB())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(500))
	}
	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"message": "user-created",
	}))
}
