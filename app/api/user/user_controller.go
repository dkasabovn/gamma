package user_api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gamma/app/api/core"
	"gamma/app/api/models/dto"
	"gamma/app/datastore/objectstore"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/argon"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (a *UserAPI) getUserController(c echo.Context) error {
	user, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	userEvents, err := a.srvc.GetUserEvents(c.Request().Context(), int(user.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	userOrgs, err := a.srvc.GetUserOrganizations(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"user":          user,
		"events":        userEvents,
		"organizations": dto.ConvertUserOrganizationRows(userOrgs),
	}))
}

func (a *UserAPI) putUserController(c echo.Context) error {
	prevUser, err := core.ExtractUser(c)
	if err != nil {
		log.Errorf("could not get user: %v", err)
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	var newUser userRepo.UpdateUserParams
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	newUser.Uuid = prevUser.Uuid
	if newUser.Email == "" {
		newUser.Email = prevUser.Email
	}
	if newUser.PasswordHash == "" {
		newUser.PasswordHash = prevUser.PasswordHash
	} else {
		hash, err := argon.PasswordToHash(newUser.PasswordHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
		}
		newUser.PasswordHash = hash
	}
	if newUser.PhoneNumber == "" {
		newUser.PhoneNumber = prevUser.PhoneNumber
	}
	if newUser.FirstName == "" {
		newUser.FirstName = prevUser.FirstName
	}
	if newUser.LastName == "" {
		newUser.LastName = prevUser.LastName
	}
	if newUser.ImageUrl == "" {
		newUser.ImageUrl = prevUser.ImageUrl
	}

	err = a.srvc.UpdateUser(c.Request().Context(), &newUser)
	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"newUser": newUser,
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
		"organizations": dto.ConvertUserOrganizationRows(orgs),
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

	inviteUser, err := a.srvc.GetOrgUser(c.Request().Context(), inviteParams.UserUuid, c.Param("org_uuid"))
	if err != nil {
		log.Errorf("Could not fetch org user: %v", err)
		return c.JSON(http.StatusInternalServerError, http.StatusInternalServerError)
	}

	inviteUuid := uuid.NewString()

	if err := a.srvc.CreateInvite(c.Request().Context(), &userRepo.InsertInviteParams{
		ExpirationDate: inviteParams.ExpirationDate,
		Capacity:       int32(inviteParams.Capacity),
		Uuid:           inviteUuid,
		OrgUserFk:      inviteUser.ID_2,
		EntityUuid:     inviteParams.EntityUuid,
		EntityType:     int32(bo.EVENT),
	}); err != nil {
		log.Errorf("could not create invite: %v", err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"invite_uuid": inviteUuid,
	}))
}

func (a *UserAPI) getOrgUserInvitesController(c echo.Context) error {
	org_user, err := core.ExtractOrguser(c, c.Param("org_uuid"))
	if err != nil {
		log.Errorf("user is not an org user")
		return c.JSON(http.StatusUnauthorized, core.ApiError(http.StatusUnauthorized))
	}

	var inviteGet dto.InviteGetEntity
	if err := c.Bind(&inviteGet); err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	invites, err := a.srvc.GetOrgUserInvites(c.Request().Context(), &userRepo.GetOrgUserInvitesParams{
		OrgUserFk:  org_user.ID_2,
		EntityUuid: inviteGet.EntityUuid,
	})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}
	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"invites": dto.ConvertInvites(invites),
	}))
}

func (a *UserAPI) getInviteController(c echo.Context) error {
	var inviteGet dto.InviteGet
	if err := c.Bind(&inviteGet); err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	invite, err := a.srvc.GetInvite(c.Request().Context(), inviteGet.InviteUuid)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	if invite.EntityType == int32(bo.EVENT) {
		event, err := a.srvc.GetEvent(c.Request().Context(), invite.EntityUuid)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
		}
		return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
			"invite": dto.ConvertInvite(invite),
			"entity": dto.ConvertOrgEvent(event),
		}))
	} else if invite.EntityType == int32(bo.ORGANIZATION) {
		organization, err := a.srvc.GetOrganization(c.Request().Context(), invite.EntityUuid)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
		}
		return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
			"invite": dto.ConvertInvite(invite),
			"entity": dto.ConvertOrganization(organization),
		}))
	}

	return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
}
