package core

import (
	"errors"

	"gamma/app/system/auth/ecJwt"

	"github.com/labstack/echo/v4"
)

func JwtParserFunction(auth string, c echo.Context) (interface{}, error) {
	token, valid := ecJwt.ECDSAVerify(auth)
	if !valid {
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
