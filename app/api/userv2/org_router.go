package user

import (
	"gamma/app/api/auth/ecJwt"
	"gamma/app/api/core"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func orgRoutes(e *echo.Echo) {
	grp := e.Group("/orgs")
	grp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	grp.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

}

func getOrgRouter(g *echo.Group) {
	g.GET("/org/:org_id", getOrgController)
}

func getOrgMemberRouter(g *echo.Group) {
	g.GET("/members/:org_id", getOrgMembersController)
}

func updateOrgRouter(g *echo.Group) {
	g.PUT("/:org_id", func(c echo.Context) error {
		return nil
	})
}
