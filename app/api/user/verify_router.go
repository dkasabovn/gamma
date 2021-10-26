package user

import (
	"fmt"
	"gamma/app/system/auth"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func updateTokens(id primitive.ObjectID , c echo.Context) error {
	accessToken, accessExp, err := auth.GenerateAccessToken(id)
	if err != nil {
		return err
	}
	refreshToken, refreshExp, err := auth.GenerateRefreshToken(id)
	if err != nil {
		return err
	}

	c.SetCookie(auth.UserCookie(id, accessExp))
	c.SetCookie(auth.TokenCookie(auth.AccessName, accessToken, accessExp))
	c.SetCookie(auth.TokenCookie(auth.RefreshName, refreshToken, refreshExp))
	return nil
}

func middleTokenUpdate(next echo.HandlerFunc) echo.HandlerFunc {

	// middleware to update refresh tokens


	return func (c echo.Context) error {
		if c.Get(auth.ClaimID) == nil {
			fmt.Println("id not found")
			return next(c)
		}

		u := c.Get(auth.ClaimID).(*jwt.Token)
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
					err = updateTokens(claims.ID, c)
					if err != nil {
						return next(c)
					}
				}
				
			}

		}
		return next(c)
	}
}


func OpenRoutes(e *echo.Echo) {
	
	open := e.Group("/api")
	
	{
		open.POST("/user", createUser)
		open.POST("/login", login)
	}
	
}

func JwtRoutes(e *echo.Echo) {
	authRequired := e.Group("/auth")
	authRequired.Use(middleware.JWTWithConfig(auth.CustomJwtConfig))
	authRequired.Use(middleTokenUpdate)
	
	{
		authRequired.GET("/users", getUsers)
		authRequired.GET("", getUser)
		authRequired.POST("", updateUser)
	}

}
