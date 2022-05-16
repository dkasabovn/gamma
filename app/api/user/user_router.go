package user_api

import (
	"gamma/app/api/core"
	"gamma/app/system/auth/ecJwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *UserAPI) addUserRoutes() {
	authRequired := a.echo.Group("/api")
	authRequired.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

	{
		a.getUserRouter(authRequired)
	}

}

func (a *UserAPI) getUserRouter(g *echo.Group) {
	g.GET("/user", a.getUserController)
}
