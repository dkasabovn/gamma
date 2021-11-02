package event

import "github.com/labstack/echo/v4"

// setUpInternalRouter sets up the internal router.
func setUpInternalRouter(e *echo.Echo) {
	g := e.Group("/user")
	{
		internalCreateUser(g)
	}
}

// TODO: Validate this endpoint with ssh keys and only allow Gabe service to talk to it
func internalCreateUser(g *echo.Group) {
	g.POST("/", func(c echo.Context) error {
		return nil
	})
}
