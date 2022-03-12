package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/auth"
	"gamma/app/services/user"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func signUpController(c echo.Context) error {
	var rawSignUp auth.UserSignup
	if err := c.Bind(&rawSignUp); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	tokens, err := user.GetUserService().CreateUser(c.Request().Context(), rawSignUp.Email, rawSignUp.RawPassword, rawSignUp.FirstName, rawSignUp.LastName)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"refresh_token": tokens.RefreshToken,
		"bearer_token":  tokens.BearerToken,
	}))
}

func signInController(c echo.Context) error {
	// TODOF: handle password reset / clues maybe
	var rawSignIn auth.UserSignIn
	if err := c.Bind(&rawSignIn); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	tokens, err := user.GetUserService().SignInUser(c.Request().Context(), rawSignIn.Email, rawSignIn.RawPassword)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	if tokens == nil {
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"refresh_token": tokens.RefreshToken,
		"bearer_token":  tokens.BearerToken,
	}))
}

func refreshTokenController(c echo.Context) error {
	var tokens ecJwt.GammaJwt
	if err := c.Bind(&tokens); err != nil {
		c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	_, refreshValid := ecJwt.ECDSAVerify(tokens.RefreshToken)
	if !refreshValid {
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	token, _ := ecJwt.ECDSAVerify(tokens.BearerToken)
	claims := token.Claims.(ecJwt.GammaClaims)
	accessToken, refreshToken := ecJwt.ECDSASign(&claims)

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"refresh_token": refreshToken,
		"bearer_token":  accessToken,
	}))
}
