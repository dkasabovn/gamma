package user

import (
	"gamma/app/api/auth/ecJwt"
	"gamma/app/api/core"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func userRoutes(e *echo.Echo) {
	grp := e.Group("/users")
	grp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	grp.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

	getSelfRouter(grp)
}

func getSelfRouter(g *echo.Group) {
	g.GET("/me", getSelfController)
}

// TODO: Grant
func updateSelfRouter(g *echo.Group) {
	g.PUT("/me", func(c echo.Context) error { return nil })
}