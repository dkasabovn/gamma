package event

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func setUpInfoGroup(e *echo.Echo) {
	g := e.Group("/events")
	getEvents(g)
}

func getEvents(g *echo.Group) {
	g.GET("/location", func(c echo.Context) error {
		return c.String(http.StatusOK, "howdy")
	})
}

func getUserEvents(g *echo.Group) {
	g.GET("/me", func(c echo.Context) error {
		return nil
	})
}
