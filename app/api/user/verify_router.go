package user

import (
	auth "gamma/app/api/core"

	"github.com/labstack/echo/v4"
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
	authRequired.Use(auth.JwtMiddle)
	
	{
		authRequired.GET("/users", getUsers)
		authRequired.GET("", getUser)
		authRequired.POST("", updateUser)
	}

}
