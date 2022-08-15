package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/services/user"
	"gamma/app/system/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func createInviteController(c echo.Context) error {
	var inviteCreateDto dto.InviteCreate
	if err := c.Bind(&inviteCreateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	orgUser, err := core.ExtractOrguser(c, inviteCreateDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	if err := user.GetUserService().CreateInvite(c.Request().Context(), orgUser, &inviteCreateDto); err != nil {
		log.Errorf("%v", err)
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}

func getInviteController(c echo.Context) error {
	var inviteGetDto dto.InviteGet
	if err := c.Bind(&inviteGetDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	// TODO: Daniel finish
	return nil
}
