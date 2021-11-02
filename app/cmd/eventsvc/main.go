package main

import (
	"gamma/app/api/event"
	"gamma/app/system"

	"github.com/labstack/echo/v4"
)

func main() {
	system.Initialize()
	e := echo.New()
	event.SetUpInfoGroup(e)
	e.Logger.Fatal(e.Start(":8080"))
	// gclaims := &ecJwt.GammaClaims{
	// 	Email: "goober",
	// 	Uuid:  primitive.ObjectID{},
	// }
	// accessToken, refreshToken := ecJwt.ECDSASign(gclaims)
	// fmt.Printf("%s\n\n%s\n", accessToken, refreshToken)
}
