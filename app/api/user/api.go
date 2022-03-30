package user_api

import (
	"gamma/app/services/iface"
	"gamma/app/services/user"

	"github.com/labstack/echo/v4"
)


var (
	apiInstance *UserAPI
)
type UserAPI struct {
	echo *echo.Echo
	port string
	srvc iface.UserService
}

func StartAPI(port string) {

	apiInstance = &UserAPI{
		echo : echo.New(),
		port : port,
		srvc: user.GetUserService(),
	}

	apiInstance.addOpenRoutes()
	apiInstance.addUserRoutes()

	apiInstance.echo.Logger.Fatal(apiInstance.echo.Start(port))
}
