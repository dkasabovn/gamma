package core

import (
	"errors"

	"gamma/app/api/auth/ecJwt"
	"gamma/app/domain/bo"

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

func GetTokens(c echo.Context, user *bo.PartialUser) (string, string) {
	claims := &ecJwt.GammaClaims{
		UUID:     user.UUID,
		Username: user.Username,
	}
	accessToken, refreshToken := ecJwt.ECDSASign(claims)
	return accessToken, refreshToken
}
