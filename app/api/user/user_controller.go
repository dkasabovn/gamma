package user_api

import (
	"fmt"
	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/datastore/objectstore"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (a *UserAPI) getUserController(c echo.Context) error {
	// TODO (add the left join method)
	user, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	userEvents, err := a.srvc.GetUserEvents(c.Request().Context(), int(user.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"user":   user,
		"events": userEvents,
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

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"organizations": dto.ConvertOrganizations(orgs),
	}))
}

func (a *UserAPI) getEventsController(c echo.Context) error {
	user, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}
	events, err := a.srvc.GetEvents(c.Request().Context(), int(user.ID))
	if err != nil {
		log.Errorf("could not get any events")
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}
	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"events": dto.ConvertEvents(events),
	}))
}

func (a *UserAPI) getEventsByOrgController(c echo.Context) error {
	_, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	events, err := a.srvc.GetOrganizationEvents(c.Request().Context(), c.Param("org_uuid"))
	if err != nil {
		log.Errorf("could not get organization events: %v", events)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"events": dto.ConvertOrgEvents(events),
	}))
}

func (a *UserAPI) postCreateEventController(c echo.Context) error {
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

	var event dto.ReqEvent
	if err := c.Bind(&event); err != nil {
		log.Errorf("could not bind event to ResEvent: %v", err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	file, err := c.FormFile("event_image")
	if err != nil {
		log.Errorf("could not open image: %v", err)
	}
	if file.Header.Get("Content-Type") != "image/webp" || file.Size > 1000000 {
		log.Errorf("Image is not properly formatted with type: %s and size: %d", file.Header.Get("Content-Type"), file.Size)
	}

	src, err := file.Open()
	if err != nil {
		log.Errorf("could not open image: %v", err)
		return nil
	}
	defer src.Close()

	imageUuid := uuid.NewString()

	image_data, err := ioutil.ReadAll(src)
	if err != nil {
		log.Errorf("Could not read all image_data: %v", err)
		return nil
	}

	image_path, err := a.store.Put(c.Request().Context(), fmt.Sprintf("/events/%s.webp", imageUuid), objectstore.Object{
		Data: image_data,
	})
	if err != nil {
		log.Errorf("Could not put image in objectstore: %v", err)
		return nil
	}

	eventUuid := uuid.NewString()

	log.Infof("Image Path: %v", image_path)

	if err := a.srvc.CreateEvent(c.Request().Context(), &userRepo.InsertEventParams{
		EventName:        event.EventName,
		EventDate:        event.EventDate,
		EventLocation:    event.EventLocation,
		EventDescription: event.EventDescription,
		Uuid:             eventUuid,
		EventImageUrl:    image_path,
		OrganizationFk:   org_user.OrganizationFk,
	}); err != nil {
		log.Errorf("Could not create event in db: %v", err)
		return nil
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"uuid": eventUuid,
	})
}

func (a *UserAPI) postEventInviteLinkController(c echo.Context) error {
	org_user, err := core.ExtractOrguser(c, c.Param("org_uuid"))
	if err != nil {
		log.Errorf("user is not an org user")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	policy_num := bo.PolicyNumber(org_user.PoliciesNum)
	if !policy_num.Can(bo.MODIFY_EVENTS) && !policy_num.Is(bo.SUPER_ADMIN) {
		log.Errorf("user is not authorized to create events")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	var inviteParams dto.InviteCreate
	if err := c.Bind(&inviteParams); err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	invite, err := bo.NewInviteWithReqs(bo.EVENT, bo.InviteRequirements{
		OrganizationID: inviteParams.Organization,
		UserID:         inviteParams.User,
	}).Serialize()
	if err != nil {
		log.Errorf("could not serialize invite: %v", err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	inviteUuid := uuid.NewString()

	if err := a.srvc.CreateInvite(c.Request().Context(), &userRepo.InsertInviteParams{
		ExpirationDate: inviteParams.ExpirationDate,
		Capacity:       int32(inviteParams.Capacity),
		PolicyJson:     invite,
		Uuid:           inviteUuid,
		OrgUserUuid:    inviteParams.OrgUserUuid,
		OrgFk:          org_user.OrganizationFk,
	}); err != nil {
		log.Errorf("could not create invite: %v", err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"invite_uuid": inviteUuid,
	}))
}
