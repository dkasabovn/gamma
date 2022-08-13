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
		a.getEventsByOrgRouter(authRequired)
		a.createEventRouter(authRequired)
		a.postEventInviteLinkRouter(authRequired)
	}
}

// @Summary Self
// @Description Get data about self and events
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <your_token>"
// @Success 200
// @Router /api/user [get]
func (a *UserAPI) getUserRouter(g *echo.Group) {
	g.GET("/user", a.getUserController)
}

// @Summary Self Orgs
// @Description Get the organizations you are in
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <your_token>"
// @Success 200
// @Router /api/orgs [get]
func (a *UserAPI) getUserOrganizationsRouter(g *echo.Group) {
	g.GET("/orgs", a.getUserOrganizationsController)
}

// @Summary Events
// @Description Get a list of events going on
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <your_token>"
// @Success 200
// @Router /api/events [get]
func (a *UserAPI) getEventsRouter(g *echo.Group) {
	g.GET("/events", a.getEventsController)
}

// @Summary Create Event
// @Description Create an event for a particular org
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <your_token>"
// @Param org_uuid path string true "Org uuid"
// @Param event_data body dto.ReqEvent true "Also needs 'event_image' which is a file"
// @Success 200
// @Router /api/events/{org_uuid} [post]
func (a *UserAPI) createEventRouter(g *echo.Group) {
	g.POST("/events/:org_uuid", a.postCreateEventController)
}

// @Summary Org Events
// @Description Get events for a particular organization
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <your_token>"
// @Param org_uuid path string true "Org uuid"
// @Success 200
// @Router /api/events/{org_uuid} [get]
func (a *UserAPI) getEventsByOrgRouter(g *echo.Group) {
	g.GET("/events/:org_uuid", a.getEventsByOrgController)
}

func (a *UserAPI) postEventApplicationRouter(g *echo.Group) {
	g.POST("/applications/:event_uuid", func(ctx echo.Context) error { return nil })
}

func (a *UserAPI) postEventInviteLinkRouter(g *echo.Group) {
	g.POST("/invite/events/:org_uuid", a.postEventInviteLinkController)
}

// @Summary Get Invites
// @Description Get invites for a particular event
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <your_token>"
// @Param org_uuid path string true "Org uuid"
// @Param entity_uuid query string true "Entity uuid"
// @Success 200
// @Router /api/invite/events/{org_uuid} [get]
func (a *UserAPI) getOrgUserInvitesRouter(g *echo.Group) {
	g.GET("/invite/events/:org_uuid", a.getOrgUserInvitesController)
}

// @Summary Get invite
func (a *UserAPI) getInviteRouter(g *echo.Group) {
	g.GET("/invite/:invite_uuid", a.getInviteController)
}

// TODO: Put user information
