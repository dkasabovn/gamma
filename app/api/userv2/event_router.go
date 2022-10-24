package user

import (
	"gamma/app/api/auth/ecJwt"
	"gamma/app/api/core"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func eventRoutes(e *echo.Echo) {
	grp := e.Group("/events")
	grp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	grp.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         &ecJwt.GammaClaims{},
		ParseTokenFunc: core.JwtParserFunction,
	}))

	getEventsRouter(grp)
	createEventRouter(grp)
	getEventRouter(grp)
	updateEventRouter(grp)
}

// With query option org_id
func getEventsRouter(g *echo.Group) {
	g.GET("/list", getEventsController)
}

func createEventRouter(g *echo.Group) {
	g.POST("/new", createEventController)
}

func checkRouter(g *echo.Group) {
	g.POST("/validate", checkController)
}

// TODO: Grant
func getEventRouter(g *echo.Group) {
	g.GET("/:event_uuid", getEventController)
}

// TODO: Grant
func updateEventRouter(g *echo.Group) {
	g.PUT("/:event_uuid", updateEventController)
}

// TODO: Grant
func deleteEventRouter(g *echo.Group) {
	g.DELETE("/:event_uuid", func(c echo.Context) error {
		return nil
	})
}
