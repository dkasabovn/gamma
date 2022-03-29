package user_api

import (
	"sync"

	"github.com/labstack/echo/v4"
)


var (
	apiOnce sync.Once
	apiInstance *UserAPI
)
type UserAPI struct {
	echo *echo.Echo
}

func API() *UserAPI{
	apiOnce.Do(func() {

		e := echo.New()
		apiInstance = &UserAPI{
			echo : e,
		}
	
		apiInstance.addOpenRoutes()
		apiInstance.addUserRoutes()

	})
	
	return apiInstance
}

func (a *UserAPI) Start(port string) {
		a.echo.Logger.Fatal(a.echo.Start(port))
} 

func (a *UserAPI) E() *echo.Echo{
	return a.echo
}
