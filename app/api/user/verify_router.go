package user

import (
	"gamma/app/api/core"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddOpenRoutes(e *echo.Echo) {

	group := e.Group("/auth")

	{
		signupRouter(group)
		loginRouter(group)
	}

}

func JwtRoutes(e *echo.Echo) {
	authRequired := e.Group("/api")
	authRequired.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

}

func signupRouter(g *echo.Group) {
	g.POST("/signup", func(c echo.Context) error {
		core.AddTokens(c, bo.User{
			Uuid:  "test",
			Email: "dkn",
		})
		return c.JSON(http.StatusOK, map[string]string{
			"howdy": "partner",
		})
	})
}

func loginRouter(g *echo.Group) {
	g.POST("/signin", func(c echo.Context) error { return nil })
}
