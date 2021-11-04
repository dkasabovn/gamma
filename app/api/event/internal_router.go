package event

import (
	"gamma/app/api/core"
	"gamma/app/system/auth/ecJwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// setUpInternalRouter sets up the internal router.
func SetUpInternalRouter(e *echo.Echo) {
	g := e.Group("/internal")
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.InternalJwtParserFunction,
	}))
	{
		internalCreateUser(g)
	}
}

// TODO: Validate this endpoint with ssh keys and only allow Gabe service to talk to it
func internalCreateUser(g *echo.Group) {
	g.POST("/", createUser)
}
