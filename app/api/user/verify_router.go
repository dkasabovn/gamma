package user

import (
	"gamma/app/api/core"
	"gamma/app/system/auth/ecJwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func OpenRoutes(e *echo.Echo) {

	open := e.Group("/auth")

	{
		open.POST("/user", createUser)
		open.POST("/login", login)
	}

}

func JwtRoutes(e *echo.Echo) {
	authRequired := e.Group("/api")
	authRequired.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

	{
		authRequired.GET("/users", getUsers)
		authRequired.GET("", getUser)
		authRequired.POST("", updateUser)
	}

}
