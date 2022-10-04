package user

import (
	"net/http"

	"gamma/app/api/auth/ecJwt"
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"gamma/app/system/log"

	"github.com/labstack/echo/v4"
)

func logInController(c echo.Context) error {
	var logInDto dto.UserSignIn
	if err := c.Bind(&logInDto); err != nil {
		log.Errorf("%v", err)
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	userData, err := user.GetUserService().SignInUser(c.Request().Context(), logInDto.Email, logInDto.RawPassword)
	if err != nil {
		log.Errorf("could not sign in user: %v", err)
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	accessToken, refreshToken := core.GetTokens(c, userData)

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"bearer_token": accessToken,
	}))
}

func signUpController(c echo.Context) error {
	var signUpDto dto.UserSignUp
	if err := c.Bind(&signUpDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	userData, err := user.GetUserService().SignUpUser(c.Request().Context(), &signUpDto)
	if err != nil {
		log.Errorf("could not create user: %v", err)
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	accessToken, refreshToken := core.GetTokens(c, userData)

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"bearer_token": accessToken,
	}))
}

func refreshController(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Errorf("could not get refresh token")
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	_, refreshValid := ecJwt.ECDSAVerify(refreshToken.Value)
	if !refreshValid {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	token, _ := ecJwt.ECDSAVerify(refreshToken.Value)
	claims := token.Claims.(*ecJwt.GammaClaims)

	user, err := user.GetUserService().GetUser(c.Request().Context(), claims.UUID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	accessToken, refreshedToken := core.GetTokens(c, &bo.PartialUser{
		UUID:     user.ID,
		Username: user.Username,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshedToken,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"bearer_token": accessToken,
	}))
}

func recoverPasswordController(c echo.Context) error {
	var resetPasswordDto dto.UserResetPasswordPreflight
	if err := c.Bind(&resetPasswordDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	if err := user.GetUserService().ResetPasswordPreflight(c.Request().Context(), &resetPasswordDto); err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}

func resetPasswordController(c echo.Context) error {
	var resetPasswordDto dto.UserResetPasswordConfirmed
	if err := c.Bind(&resetPasswordDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	if err := user.GetUserService().ResetPasswordConfirmed(c.Request().Context(), &resetPasswordDto); err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}
