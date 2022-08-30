package user

import (
	"fmt"
	"gamma/app/api/auth/argon"
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"gamma/app/system/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getSelfController(c echo.Context) error {
	self, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	userEvents, err := user.GetUserService().GetUserEvents(c.Request().Context(), self.ID)
	if err != nil {
		log.Errorf("could not get user events: %v", err)
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	userOrgs, err := user.GetUserService().GetUserOrganizations(c.Request().Context(), self.ID)
	if err != nil {
		log.Errorf("could nto get user organizations: %v", err)
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"user":               dto.ConvertUser(self),
		"user_events":        dto.ConvertUserEvents(userEvents),
		"user_organizations": dto.ConvertOrganizationsWithPermissions(userOrgs),
	}))
}

func updateSelfController(c echo.Context) error {
	prevUser, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}
	fmt.Print(prevUser.Email)

	var newUser dto.UserUpdate
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	if newUser.Email != "" {
		prevUser.Email = newUser.Email
	}

	if newUser.FirstName != "" {
		prevUser.FirstName = newUser.FirstName
	}

	if newUser.LastName != "" {
		prevUser.LastName = newUser.LastName
	}

	if newUser.RawPassword != "" {
		hash, err := argon.PasswordToHash(newUser.RawPassword)
		if err != nil {
			log.Infof("%v", err)
			return err
		}
		prevUser.PasswordHash = hash
	}

	if newUser.UserName != "" {
		prevUser.Username = newUser.UserName
	}

	if newUser.ImageUrl != "" {
		prevUser.ImageUrl = newUser.ImageUrl
	}

	err = user.GetUserService().UpdateUser(c.Request().Context(), prevUser)
	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"newUser": prevUser,
	}))

}
