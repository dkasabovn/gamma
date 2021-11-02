package event

import (
	"gamma/app/api/core"

	"github.com/labstack/echo/v4"
)

func SetUpInfoGroup(e *echo.Echo) {
	g := e.Group("/events")
	g.Use(core.JwtMiddle)
	getApplications(g)
	getUserEvents(g)
	getOrgs(g)
}

func getApplications(g *echo.Group) {
	g.GET("/applications", GetEventApplications)
}

func getUserEvents(g *echo.Group) {
	g.GET("", GetAttendingEvents)
}

func getOrgs(g *echo.Group) {
	g.GET("/organizations", GetBootstrapData)
}
