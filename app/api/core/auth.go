package refresh

import (
	"fmt"
	auth "gamma/app/system/auth"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

/*
Given the users previous claims updates jwt tokens
*/
func UpdateTokens(claim auth.UserClaims, c echo.Context) error {
	accessToken, accessExp, err := auth.GenerateAccessToken(claim)
	if err != nil {
		fmt.Println("error while making access")
		return err
	}
	refreshToken, refreshExp, err := auth.GenerateRefreshToken(claim)
	if err != nil {
		fmt.Println("error while making refresh")
		return err
	}
	
	
	c.SetCookie(auth.TokenCookie(auth.AccessName, accessToken, accessExp))
	c.SetCookie(auth.TokenCookie(auth.RefreshName, refreshToken, refreshExp))
	
	return nil
}


/*
Middleware method to validate JWT and update tokens if needed
return status.unauthorized if refresh tokens are bad, or if tokens could no be updated
*/
func JwtMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {


		refreshCookie, err := c.Cookie(auth.RefreshName);
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Log in first");
		}
		accessCookie, err := c.Cookie(auth.AccessName);
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Log in first");
		}

		
		refreshToken, err := jwt.ParseWithClaims(refreshCookie.Value, &auth.UserClaims{}, auth.KeyFunc);
		if err != nil || refreshToken == nil {
			fmt.Print(err)
			return  c.JSON(http.StatusUnauthorized, "Expired Refresh");
		}

		claims := refreshToken.Claims.(*auth.UserClaims)
		accessToken, err := jwt.ParseWithClaims(accessCookie.Value, &auth.UserClaims{}, auth.KeyFunc);
		if err != nil || accessToken == nil{
			if err = UpdateTokens(*claims, c); err != nil {
				return  c.JSON(http.StatusUnauthorized, "Expired Access and could not update tokens, Login in again ");
			}
		} else if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15 * time.Minute {
			if err = UpdateTokens(*claims, c); err != nil {
				return c.JSON(http.StatusBadRequest, "Bad Claims")
			}
		}
		return next(c)

	}
}

// func MiddleTokenUpdate(next echo.HandlerFunc) echo.HandlerFunc {

// 	// middleware to update refresh tokens


// 	return func (c echo.Context) error {
// 		for _, cookie := range c.Cookies() {
// 		fmt.Printf(">%s<\n", cookie.Name)
// 		fmt.Println(cookie.Value)
// 	}
// 		if c.Get("access-token") == nil {
// 			fmt.Println("access not found")
// 			return next(c)
// 		}
// 		fmt.Println("Found")


// 		u := c.Get(auth.RefreshName).(*jwt.Token)
// 		claims := u.Claims.(*auth.UserClaims)

// 		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15 * time.Minute {

// 			refreshCookie, err := c.Cookie(auth.RefreshName)
// 			if err == nil && refreshCookie != nil {

// 				refreshTkn, err := jwt.ParseWithClaims(refreshCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
// 					return []byte(auth.GetJwtRefresh()), nil
// 				})

// 				if err != nil {
// 					c.Response().Writer.WriteHeader(http.StatusUnauthorized)
// 				}

// 				if refreshTkn != nil && refreshTkn.Valid {
// 					err = UpdateTokens(*claims, c)
// 					if err != nil {
// 						return next(c)
// 					}
// 				}
				
// 			}

// 		}
// 		return next(c)
// 	}
// }
