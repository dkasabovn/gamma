package user

import (
	"gamma/app/api/core"
	"gamma/app/system/auth/ecJwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddOpenRoutes(e *echo.Echo) {

	group := e.Group("/auth")

	{
		signupRouter(group)
		loginRouter(group)
		refreshRouter(group)
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
	g.POST("/signup", signUpController)
}

func loginRouter(g *echo.Group) {
	g.POST("/signin", signInController)
}

func refreshRouter(g *echo.Group) {
	g.GET("/refresh", refreshTokenController)
}
