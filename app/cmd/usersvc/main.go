package main

import (
	"gamma/app/api/user"
	"net/http"

	"github.com/labstack/echo/v4"
<<<<<<< HEAD
	"github.com/labstack/echo/v4/middleware"
=======
>>>>>>> 365cd9b424e878144faca027c882549217cb0266
)

func main() {

	e := echo.New()

<<<<<<< HEAD
	authRequired := e.Group("/api")

	// Configure middleware with the custom claims type
	jwtConfig := middleware.JWTConfig{
		Claims:      &auth.Claims{},
		SigningKey:  []byte(auth.GetJwtSecret()),
		TokenLookup: "cookie:access-token",
	}
	authRequired.Use(middleware.JWTWithConfig(jwtConfig))
	authRequired.GET("", authNeed)
=======
	//any request should try to update the tokens
	e.Use(user.MiddleTokenUpdate)

	// adds temp get and post routes
	user.JwtRoutes(e)
>>>>>>> 365cd9b424e878144faca027c882549217cb0266

	// temp no auth rout
	e.GET("/", noAuth)

	e.Logger.Fatal(e.Start(":8000"))

}

func noAuth(ctx echo.Context) error {
	return ctx.JSON(http.StatusAccepted, "No auth needed")
<<<<<<< HEAD
}

func authNeed(ctx echo.Context) error {
	userCookie, _ := ctx.Cookie("user")
	return ctx.JSON(http.StatusAccepted, fmt.Sprintf("You have been authenticated %s", userCookie.Value))
}
=======
}
>>>>>>> 365cd9b424e878144faca027c882549217cb0266
