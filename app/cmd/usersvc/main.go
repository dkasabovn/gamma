package main

import (
	"fmt"
	"gamma/app/api/user"
	"gamma/app/datastore/users"
	"gamma/app/system"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	system.Initialize()
	users.MongoDB()
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("", func(c echo.Context) error {
		fmt.Println("Printing Cookies")
		fmt.Println(c.Cookies())
		for _, cookie := range c.Cookies() {
			fmt.Println(cookie.Name)
			fmt.Println(cookie.Value)
		}
		return nil
	})

	user.JwtRoutes(e)
	user.OpenRoutes(e)

	e.Logger.Fatal(e.Start(":8000"))

}
