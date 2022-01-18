package event

import (
	"gamma/app/api/core"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetUpInfoGroup(e *echo.Echo) {
	g := e.Group("/events")
	// JWT middleware.
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))
	{
		getTest(g)
	}
}

func getTest(g *echo.Group) {
	g.GET("/test", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*ecJwt.GammaClaims)
		return c.String(http.StatusOK, claims.Email)
	})
}
