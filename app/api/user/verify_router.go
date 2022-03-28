package user

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddOpenRoutes(e *echo.Echo) {

	group := e.Group("/auth")

	group.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	{
		signupRouter(group)
		loginRouter(group)
		refreshRouter(group)
	}

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
