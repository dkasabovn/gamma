package main

import (
	"gamma/app/api/user"
	userDB "gamma/app/datastore/user"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	userDB.MongoUsers()


	user.JwtRoutes(e)
	user.OpenRoutes(e)
	
	e.Logger.Fatal(e.Start(":8000"))

}
