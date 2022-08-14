package core

import (
	"errors"
	"gamma/app/api/auth/ecJwt"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/services/user"
	"io/ioutil"
	"net/http"

	"gamma/app/system/log"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ApiResponse struct {
	Data         map[string]interface{} `json:"data"`
	Error        int                    `json:"error_code"`
	ErrorMessage *string                `json:"error_message,omitempty"`
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

func JSONApiError(c echo.Context, errorCode int) error {
	return c.JSON(errorCode, ApiError(errorCode))
}

func FormImage(c echo.Context, data []byte, name string) error {
	file, err := c.FormFile(name)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	if file.Header.Get("Content-Type") != "image/webp" || file.Size > 10000000 {
		log.Errorf("Image with content type %s and size %d is invalid", file.Header.Get("Coontent-Type"), file.Size)
		return errors.New("image isn't the right format")
	}

	src, err := file.Open()
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	defer src.Close()

	imageData, err := ioutil.ReadAll(src)
	if err != nil {
		log.Errorf("Could not read all image_data: %v", err)
		return err
	}

	data = imageData
	return nil
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
	user, err := user.GetUserService().GetUser(c.Request().Context(), claims.UUID)
	if err != nil {
		log.Errorf("Could not extract user from faulty token: %s", userToken.Raw)
		return nil, err
	}
	return user, nil
}

func ExtractOrguser(c echo.Context, orgUuidString string) (*userRepo.GetOrgUserRow, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*ecJwt.GammaClaims)
	orgUuid, err := uuid.Parse(orgUuidString)
	if err != nil {
		log.Errorf("org uuid is not in uuid format: %v", err)
		return nil, err
	}
	user, err := user.GetUserService().GetOrgUser(c.Request().Context(), claims.UUID, orgUuid)
	if err != nil {
		log.Errorf("Could not extract org user from token: %s; %v", userToken.Raw, err)
		return nil, err
	}
	return user, nil
}
