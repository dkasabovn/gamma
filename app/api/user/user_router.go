package user_api

import (
	"gamma/app/api/core"
	"gamma/app/system/auth/ecJwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *UserAPI) addUserRoutes() {
	authRequired := a.echo.Group("/api")
	authRequired.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	authRequired.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

	{
		a.getUserRouter(authRequired)
		a.getEventsRouter(authRequired)
		a.getUserOrganizationsRouter(authRequired)
		a.getOrgImageUploadRouter(authRequired)
		a.getEventsByOrgRouter(authRequired)
	}

}

func (a *UserAPI) getUserRouter(g *echo.Group) {
	g.GET("/user", a.getUserController)
}

func (a *UserAPI) getUserOrganizationsRouter(g *echo.Group) {
	g.GET("/orgs", a.getUserOrganizationsController)
}

func (a *UserAPI) getEventsRouter(g *echo.Group) {
	g.GET("/events", a.getEventsController)
}

func (a *UserAPI) createEventRouter(g *echo.Group) {
	g.POST("/event/:org_uuid", a.createEventController)
}

func (a *UserAPI) getEventsByOrgRouter(g *echo.Group) {
	g.GET("/events/:org_uuid", a.getEventsByOrgController)
}

func (a *UserAPI) getEventApplicationsRouter(g *echo.Group) {
	g.GET("/applications/:event_uuid", func(ctx echo.Context) error { return nil })
}

func (a *UserAPI) postEventApplicationDecisionsRouter(g *echo.Group) {
	g.POST("/applications/:event_uuid", func(ctx echo.Context) error { return nil })
}

func (a *UserAPI) getOrgImageUploadRouter(g *echo.Group) {
	g.GET("/orgimage/:org_uuid", a.getOrgImageUploadController)
}
