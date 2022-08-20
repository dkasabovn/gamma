package user

import (
	"gamma/app/api/auth/ecJwt"
	"gamma/app/api/core"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func inviteRoutes(e *echo.Echo) {
	grp := e.Group("/invites")
	grp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	grp.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

	getInviteRouter(grp)
	createInviteRouter(grp)
}

func getInviteRouter(g *echo.Group) {
	g.GET("/invite/:invite_id", getInviteController)
}

func getSelfInvitesRouter(g *echo.Group) {
	g.GET("/me", func(c echo.Context) error {
		return nil
	})
}

func acceptInviteRouter(g *echo.Group) {
	g.GET("/accept/:invite_id", func(c echo.Context) error {
		return nil
	})
}

func createInviteRouter(g *echo.Group) {
	g.POST("/new", createInviteController)
}
