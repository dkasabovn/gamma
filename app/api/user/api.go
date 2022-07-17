package user_api

import (
	"gamma/app/datastore/objectstore"
	"gamma/app/services/iface"
	"gamma/app/services/user"

	"github.com/labstack/echo/v4"

	_ "gamma/app/cmd/user/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	apiInstance *UserAPI
)

type UserAPI struct {
	echo  *echo.Echo
	port  string
	srvc  iface.UserService
	store objectstore.Storage
}

// @title Gamma User Api
// @version 0.0
// @description The api docs
// @BasePath /

func StartAPI(port string) {

	apiInstance = &UserAPI{
		echo:  echo.New(),
		port:  port,
		srvc:  user.GetUserService(),
		store: objectstore.GetStorage(),
	}

	apiInstance.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	apiInstance.addOpenRoutes()
	apiInstance.addUserRoutes()

	apiInstance.echo.Logger.Fatal(apiInstance.echo.Start(port))
}
