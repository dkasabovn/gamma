package user_api

import (
	"net/http"

	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/ecJwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (a *UserAPI) signUpController(c echo.Context) error {
	var rawSignUp dto.UserSignup
	if err := c.Bind(&rawSignUp); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	// TODO(dk): upload the image and create the url instead of passing it in through raw signup
	tokens, err := a.srvc.CreateUser(c.Request().Context(), rawSignUp.RawPassword, rawSignUp.Email, rawSignUp.PhoneNumber, rawSignUp.FirstName, rawSignUp.LastName, rawSignUp.UserName, rawSignUp.ImageUrl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"bearer_token": tokens.BearerToken,
	}))
}

func (a *UserAPI) signInController(c echo.Context) error {
	// TODOF: handle password reset / clues maybe
	var rawSignIn dto.UserSignIn
	if err := c.Bind(&rawSignIn); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	tokens, err := a.srvc.SignInUser(c.Request().Context(), rawSignIn.Email, rawSignIn.RawPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	if tokens == nil {
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"bearer_token": tokens.BearerToken,
	}))
}

func (a *UserAPI) refreshTokenController(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	log.Infof("Cookies: %v", c.Cookies())
	if err != nil {
		log.Errorf("Could not get refresh_token")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	_, refreshValid := ecJwt.ECDSAVerify(refreshToken.Value)
	if !refreshValid {
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	token, _ := ecJwt.ECDSAVerify(refreshToken.Value)
	claims := token.Claims.(*ecJwt.GammaClaims)

	var user *bo.User
	if user, err = a.srvc.GetUser(c.Request().Context(), claims.Uuid); err != nil {
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	tokens := ecJwt.GetTokens(c.Request().Context(), claims.Uuid, user.Email, user.UserName, "https://tinyurl.com/monkeygamma")

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"bearer_token": tokens.BearerToken,
	}))
}
