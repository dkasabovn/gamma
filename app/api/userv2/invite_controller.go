package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/domain/bo"
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

	_, err := core.ExtractUser(c)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	invite, err := user.GetUserService().GetInvite(c.Request().Context(), &inviteGetDto)
	if err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	var entity any

	// TODO: Typeify this string
	switch invite.EntityType {
	case int32(bo.EVENT):
		event, err := user.GetUserService().GetEvent(c.Request().Context(), invite.EntityUuid)
		if err != nil {
			return core.JSONApiError(c, http.StatusInternalServerError)
		}
		entity = dto.ConvertEvent(event)
	case int32(bo.ORGANIZATION):
		organiztion, err := user.GetUserService().GetOrganization(c.Request().Context(), invite.EntityUuid)
		if err != nil {
			return core.JSONApiError(c, http.StatusInternalServerError)
		}
		entity = dto.ConvertOrganization(organiztion)
	default:
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"entity": entity,
		"invite": dto.ConvertInvite(invite),
	}))
}

func getSelfInvitesController(c echo.Context) error {
	self, err := core.ExtractUser(c)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	invites, err := user.GetUserService().GetInvitesForOrgUser(c.Request().Context(), self.ID)
	if err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"invites": dto.ConvertInvites(invites),
	}))
}

func acceptInviteController(c echo.Context) error {
	self, err := core.ExtractUser(c)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	var inviteGetDto dto.InviteGet
	if err := c.Bind(&inviteGetDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	if err := user.GetUserService().AcceptInvite(c.Request().Context(), self, &inviteGetDto); err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{}))
}
