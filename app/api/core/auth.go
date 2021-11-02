package core

import (
	"errors"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func JwtParserFunction(auth string, c echo.Context) (interface{}, error) {
	token, valid := ecJwt.ECDSAVerify(auth)
	if !valid {
		if refreshTokenString, err := c.Cookie("refresh_token"); err != nil {
			_, refreshValid := ecJwt.ECDSAVerify(refreshTokenString.Value)
			if refreshValid {
				claims := token.Claims.(ecJwt.GammaClaims)
				accessToken, refreshToken := ecJwt.ECDSASign(&claims)
				c.SetCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: refreshToken,
				})
				newToken, newValid := ecJwt.ECDSAVerify(accessToken)
				if newValid {
					return newToken, nil
				}
			}
		}
		return nil, errors.New("token is invalid")
	}
	return token, nil
}
