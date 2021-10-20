package main

import (
	"fmt"
	"gamma/app/system/auth"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)



func main() {

	e := echo.New()


	authRequired := e.Group("/api")

	// Configure middleware with the custom claims type
	jwtConfig := middleware.JWTConfig{
		Claims:     &auth.Claims{},
		SigningKey: []byte(auth.GetJwtSecret()),
		TokenLookup: "cookie:access-token",
	}
	authRequired.Use(middleware.JWTWithConfig(jwtConfig))
	authRequired.GET("", authNeed)

	e.GET("/", noAuth)

	e.Logger.Fatal(e.Start(":8000"))


}

func noAuth(ctx echo.Context) error {
	return ctx.JSON(http.StatusAccepted, "No auth needed")
}

func authNeed(ctx echo.Context) error {
	userCookie, _:= ctx.Cookie("user")
	return ctx.JSON(http.StatusAccepted, fmt.Sprintf("You have been authenticated %s", userCookie.Value ))
}