package user

import (
	"gamma/app/system/auth"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func updateTokens(email string , c echo.Context) error {
	accessToken, accessExp, err := auth.GenerateAccessToken(email)
	if err != nil {
		return err
	}
	refreshToken, refreshExp, err := auth.GenerateRefreshToken(email)
	if err != nil {
		return err
	}

	c.SetCookie(auth.UserCookie(email, accessExp))
	c.SetCookie(auth.TokenCookie(auth.AccessName, accessToken, accessExp))
	c.SetCookie(auth.TokenCookie(auth.RefreshName, refreshToken, refreshExp))
	return nil
}

func MiddleTokenUpdate(next echo.HandlerFunc) echo.HandlerFunc {

	// middleware to update refresh tokens


	return func (c echo.Context) error {
		if c.Get("email") == nil {
			return next(c)
		}

		u := c.Get("email").(*jwt.Token)
		claims := u.Claims.(*auth.Claims)

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15 * time.Minute {

			refreshCookie, err := c.Cookie(auth.RefreshName)
			if err == nil && refreshCookie != nil {

				refreshTkn, err := jwt.ParseWithClaims(refreshCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(auth.GetJwtRefresh()), nil
				})

				if err != nil {
					c.Response().Writer.WriteHeader(http.StatusUnauthorized)
				}

				if refreshTkn != nil && refreshTkn.Valid {
					err = updateTokens(claims.Email, c)
					if err != nil {
						return next(c)
					}
				}
				
			}

		}
		return next(c)
	}
}


func JwtRoutes(e *echo.Echo) {
	authRequired := e.Group("/api")
	authRequired.Use(middleware.JWTWithConfig(auth.CustomJwtConfig))

	{
		user(authRequired)
	}

}

func user(group *echo.Group) {
	group.GET("/", getUser)
	group.POST("/", createUser)
}