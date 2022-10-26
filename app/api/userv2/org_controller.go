package user

import (
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func createOrgController(c echo.Context) error {
	var orgCreateDto dto.OrganizationCreate
	if err := c.Bind(&orgCreateDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}
	orgUUID := uuid.New()
	orgParams := &userRepo.InsertOrganizationParams{
		OrgName:     orgCreateDto.OrgName,
		City:        orgCreateDto.City,
		OrgImageUrl: orgCreateDto.OrgImageUrl,
		ID:          orgUUID,
	}
	err := user.GetUserService().CreateOrganization(c.Request().Context(), orgParams)
	if err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}
	userUUID, err := uuid.Parse(orgCreateDto.AdminId)
	if err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}
	orgUserParams := &userRepo.InsertOrgUserParams{
		PoliciesNum:    bo.ADMIN,
		UserFk:         userUUID,
		OrganizationFk: orgUUID,
	}
	err = user.GetUserService().CreateOrgUser(c.Request().Context(), orgUserParams)
	if err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"organization": "created",
		"admin":        "assigned",
	}))
}

func getOrgController(c echo.Context) error {
	var orgGetDto dto.OrganizationGet
	if err := c.Bind(&orgGetDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	orgUUID, err := uuid.Parse(orgGetDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	org, err := user.GetUserService().GetOrganization(c.Request().Context(), orgUUID)
	if err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"organization": dto.ConvertOrganization(org),
	}))
}

func getOrgMembersController(c echo.Context) error {
	var orgGetDto dto.OrganizationGet
	if err := c.Bind(&orgGetDto); err != nil {
		return core.JSONApiError(c, http.StatusBadRequest)
	}

	orgUser, err := core.ExtractOrguser(c, orgGetDto.OrganizationID)
	if err != nil {
		return core.JSONApiError(c, http.StatusUnauthorized)
	}

	orgMembers, err := user.GetUserService().GetOrganizationUsers(c.Request().Context(), orgUser.OrganizationFk)
	if err != nil {
		return core.JSONApiError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"members": orgMembers,
	}))
}
