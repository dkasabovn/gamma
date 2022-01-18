package core

import (
	"errors"
	"fmt"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func JwtParserFunction(auth string, c echo.Context) (interface{}, error) {
	token, valid := ecJwt.ECDSAVerify(auth)
	if !valid {
		if refreshTokenString, err := c.Cookie("refresh-token"); err == nil {
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

func InternalJwtParserFunction(auth string, c echo.Context) (interface{}, error) {
	token, valid := ecJwt.ECDSAVerify(auth)
	if !valid {
		return nil, errors.New("token is invalid")
	}
	gammaClaims := token.Claims.(ecJwt.GammaClaims)
	if gammaClaims.Audience != "usersvc.gamma" {
		return nil, errors.New("invalid access to internal function")
	}
	return token, nil
}

func AddTokens(c echo.Context, user bo.User) {
	claims := &ecJwt.GammaClaims{
		Email: user.Email,
		Uuid:  user.Uuid,
	}

	accessToken, refreshToken := ecJwt.ECDSASign(claims)
	c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	c.SetCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})
}
