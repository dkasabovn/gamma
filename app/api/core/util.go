package core

import (
	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/ecJwt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ApiResponse struct {
	Data  map[string]interface{} `json:"data"`
	Error int                    `json:"error_code"`
}

// ApiConverter converts a map to a response object
func ApiConverter(data map[string]interface{}, errorCode int) *ApiResponse {
	return &ApiResponse{
		Data:  data,
		Error: errorCode,
	}
}

func ApiSuccess(data map[string]interface{}) *ApiResponse {
	return &ApiResponse{
		Data:  data,
		Error: 0,
	}
}

func ApiError(errorCode int) *ApiResponse {
	return &ApiResponse{
		Data:  nil,
		Error: errorCode,
	}
}

func GetGammaClaims(c echo.Context) *ecJwt.GammaClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*ecJwt.GammaClaims)
	return claims
}

func ExtractUser(c echo.Context) (*bo.User, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*ecJwt.GammaClaims)
	user, err := pg.GetUserRepo().GetUser(c.Request().Context(), claims.Uuid)
	if err != nil {
		log.Errorf("Could not extract user from faulty token: %s", userToken.Raw)
		return nil, err
	}
	return user, nil
}
