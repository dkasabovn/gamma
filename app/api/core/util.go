package core

import (
	userRepo "gamma/app/datastore/pg"
	"gamma/app/services/user"
	"gamma/app/system/auth/ecJwt"
	"net/http"

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

func GetCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func ExtractUser(c echo.Context) (*userRepo.User, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*ecJwt.GammaClaims)
	user, err := user.GetUserService().GetUser(c.Request().Context(), claims.Uuid)
	if err != nil {
		log.Errorf("Could not extract user from faulty token: %s", userToken.Raw)
		return nil, err
	}
	return user, nil
}

func ExtractOrguser(c echo.Context, org_uuid string) (*userRepo.GetUserOrgUserJoinRow, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*ecJwt.GammaClaims)
	user, err := user.GetUserService().GetUserOrgUserByUuid(c.Request().Context(), claims.Uuid, org_uuid)
	if err != nil {
		log.Errorf("Could not extract org user from token: %s", userToken.Raw)
		return nil, err
	}
	return user, nil
}
