package user_api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *UserAPI) addOpenRoutes() {

	group := a.echo.Group("/auth")

	group.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	{
		a.signupRouter(group)
		a.loginRouter(group)
		a.refreshRouter(group)
	}

}

func (a *UserAPI) signupRouter(g *echo.Group) {
	g.POST("/signup", SignUpController)
}

func (a *UserAPI) loginRouter(g *echo.Group) {
	g.POST("/signin", SignInController)
}

func (a *UserAPI) refreshRouter(g *echo.Group) {
	g.GET("/refresh", RefreshTokenController)
}
