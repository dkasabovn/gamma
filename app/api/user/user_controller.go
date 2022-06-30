package user_api

import (
	"fmt"
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/datastore/objectstore"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/ecJwt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (a *UserAPI) getUserController(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*ecJwt.GammaClaims)
	uuid := claims.Uuid

	user, err := a.srvc.GetUser(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"user": user,
	}))
}

func (a *UserAPI) getUserOrganizationsController(c echo.Context) error {
	org_user, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	orgs, err := a.srvc.GetUserOrganizations(c.Request().Context(), org_user.ID)
	if err != nil {
		log.Errorf("could not find organizations for user: %v", err)
		return c.JSON(http.StatusNotFound, core.ApiError(http.StatusNotFound))
	}

	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"organizations": dto.ConvertOrganizations(orgs),
	}))
}

func (a *UserAPI) getEventsController(c echo.Context) error {
	events, err := a.srvc.GetEvents(c.Request().Context())
	if err != nil {
		log.Errorf("could not get any events")
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}
	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"events": dto.ConvertEvents(events),
	}))
}

func (a *UserAPI) createEventController(c echo.Context) error {
	org_user, err := core.ExtractOrguser(c, c.Param("org_uuid"))
	if err != nil {
		log.Errorf("user is not an org user")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	policy_num := bo.PolicyNumber(org_user.PoliciesNum)
	if !policy_num.Can(bo.CREATE_EVENTS) && !policy_num.Is(bo.SUPER_ADMIN) {
		log.Errorf("user is not authorized to create events")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	var rawEvent dto.ResEvent
	if err := c.Bind(&rawEvent); err != nil {
		log.Errorf("improper event formatting: %v", err)
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	eventUuid := uuid.NewString()

	err = a.srvc.CreateEvent(c.Request().Context(), &userRepo.InsertEventParams{
		EventName:        rawEvent.EventName,
		EventDate:        rawEvent.EventDate,
		EventLocation:    rawEvent.EventLocation,
		EventDescription: rawEvent.EventDescription,
		Uuid:             eventUuid,
		EventImageUrl:    rawEvent.EventImageUrl,
		OrganizationFk:   org_user.OrganizationFk,
	})

	if err != nil {
		log.Errorf("could not create event: %v", err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"event": eventUuid,
	}))
}

func (a *UserAPI) getEventsByOrgController(c echo.Context) error {
	events, err := a.srvc.GetOrganizationEvents(c.Request().Context(), c.Param("org_uuid"))
	if err != nil {
		log.Errorf("could not get organization events: %v", events)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"events": dto.ConvertEvents(events),
	}))
}

func (a *UserAPI) getOrgImageUploadController(c echo.Context) error {
	user, err := core.ExtractOrguser(c, c.Param("org_uuid"))
	if err != nil {
		log.Errorf("user is not an org user")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	policy_num := bo.PolicyNumber(user.PoliciesNum)
	if !policy_num.Can(bo.CREATE_EVENTS) {
		log.Errorf("user is not authorized to create events")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	upload_url := fmt.Sprintf("orgs/%s/%s", c.Param("org_uuid"), uuid.NewString())
	url, err := objectstore.GenerateObjectUploadUrl(upload_url)
	if err != nil {
		log.Errorf("could not generate presigned url")
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}
