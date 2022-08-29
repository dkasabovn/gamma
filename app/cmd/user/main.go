package main

import (
	user "gamma/app/api/userv2"
	"log"

	"github.com/labstack/echo/v4"
)

const (
	port = ":8080"
)

func main() {
	e := echo.New()
	user.AddRoutes(e)
	log.Fatalf("%v", e.Start(port))
}
