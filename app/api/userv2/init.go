package user

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo) {
	eventRoutes(e)
	authRoutes(e)
	inviteRoutes(e)
	userRoutes(e)
}
