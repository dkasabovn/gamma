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

// @Summary Sign Up
// @Description Sign Up User
// @Accept json
// @Produce json
// @Param Details body dto.UserSignUp true "Data required to create an account"
// @Success 200
// @Router /auth/signup [post]
func (a *UserAPI) signupRouter(g *echo.Group) {
	g.POST("/signup", a.signUpController)
}

// @Summary Sign In
// @Description Sign In User
// @Accept json
// @Produce json
// @Param Details body dto.UserSignIn true "Username and Password"
// @Success 200
// @Router /auth/signin [post]
func (a *UserAPI) loginRouter(g *echo.Group) {
	g.POST("/signin", a.signInController)
}

// @Summary Refresh Tokens
// @Description Send in a valid refresh token (http only cookie) and get a new set of tokens
// @Accept json
// @Produce json
// @Success 200
// @Router /auth/refresh [get]
func (a *UserAPI) refreshRouter(g *echo.Group) {
	g.GET("/refresh", a.refreshTokenController)
}
