package user

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func authRoutes(e *echo.Echo) {
	grp := e.Group("/auth")
	grp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	logInRouter(grp)
	signUpRouter(grp)
	refreshRouter(grp)
}

func logInRouter(g *echo.Group) {
	g.POST("/login", logInController)
}

func signUpRouter(g *echo.Group) {
	g.POST("/signup", signUpController)
}

func refreshRouter(g *echo.Group) {
	g.GET("/refresh", refreshController)
}

func recoverPasswordRouter(g *echo.Group) {
	g.GET("/recover", func(c echo.Context) error { return nil })
}
