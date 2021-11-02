package core

import (
	"gamma/app/system/auth/ecJwt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
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
		Data: data,
	}
}

func ApiError(errorCode int) *ApiResponse {
	return &ApiResponse{
		Error: errorCode,
	}
}

func GetGammaClaims(c echo.Context) *ecJwt.GammaClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*ecJwt.GammaClaims)
	return claims
}
