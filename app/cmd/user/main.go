package main

import (
	user "gamma/app/api/userv2"
	"gamma/app/system"
	"log"

	"github.com/labstack/echo/v4"
)

const (
	port = ":8080"
)

func main() {
	system.Initialize()
	e := echo.New()
	user.AddRoutes(e)
	log.Fatalf("%v", e.Start(port))
}
