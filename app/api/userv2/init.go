package user

import (
	"gamma/app/system"
	"gamma/app/system/stability"

	"github.com/labstack/echo/v4"
)

func init() {
	system.Initialize()
	stability.LoadDependencies(stability.UserSvc())
}

func AddRoutes(e *echo.Echo) {
	eventRoutes(e)
	authRoutes(e)
	inviteRoutes(e)
	userRoutes(e)
}
