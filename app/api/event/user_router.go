package event

import "github.com/labstack/echo/v4"

func setUpUserGroup(e *echo.Echo) {
	g := e.Group("/user")
	{
		createUser(g)
		acceptOrgInvite(g)
		createApplicationForEvent(g)
		deleteApplicationForEvent(g)
	}
}

func createUser(g *echo.Group) {
	g.POST("/", func(c echo.Context) error {
		return nil
	})
}
func acceptOrgInvite(g *echo.Group) {
	g.GET("/invite/:inviteId", func(c echo.Context) error {
		return nil
	})
}

func createApplicationForEvent(g *echo.Group) {
	g.GET("/apply/:eventId", func(c echo.Context) error {
		return nil
	})
}

func deleteApplicationForEvent(g *echo.Group) {
	g.DELETE("/apply/:eventId", func(c echo.Context) error {
		return nil
	})
}
