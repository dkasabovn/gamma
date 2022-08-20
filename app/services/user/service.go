package user

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gamma/app/api/auth/argon"
	"gamma/app/api/models/dto"
	"gamma/app/datastore/objectstore"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/iface"
	"gamma/app/system/log"

	"github.com/google/uuid"
)

const (
	defaultImageUrl = "users/default.webp"
)

var (
	userOnce            sync.Once
	userServiceInstance iface.UserService
)

type userService struct {
	userRepo *userRepo.Queries
	// redis    redis.Redis
	storage objectstore.Storage
}

func GetUserService() iface.UserService {
	userOnce.Do(func() {
		userServiceInstance = &userService{
			userRepo: userRepo.New(userRepo.RwInstance()),
			// redis:    redis.GetRedis(),
			storage: objectstore.GetStorage(),
		}
	})
	return userServiceInstance
}

func (u *userService) SignInUser(ctx context.Context, email, password string) (*bo.PartialUser, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Infof("%v", err)
		return nil, err
	}

	valid, err := argon.PasswordIsMatch(password, user.PasswordHash)
	if err != nil || !valid {
		log.Infof("%v", err)
		return nil, fmt.Errorf("incorrect password: %v", err)
	}

	return &bo.PartialUser{
		UUID:     user.ID,
		Username: user.Username,
	}, nil
}

func (u *userService) SignUpUser(ctx context.Context, signUpParams *dto.UserSignUp) (*bo.PartialUser, error) {
	// TODO: user validation
	newUuid := uuid.New()

	hash, err := argon.PasswordToHash(signUpParams.RawPassword)
	if err != nil {
		log.Infof("%v", err)
		return nil, err
	}

	if err := u.userRepo.InsertUser(ctx, &userRepo.InsertUserParams{
		ID:           newUuid,
		Email:        signUpParams.Email,
		PasswordHash: hash,
		PhoneNumber:  signUpParams.PhoneNumber,
		FirstName:    signUpParams.FirstName,
		LastName:     signUpParams.LastName,
		ImageUrl:     defaultImageUrl,
		Username:     signUpParams.UserName,
		Validated:    false,
	}); err != nil {
		log.Infof("%v", err)
		return nil, err
	}

	return &bo.PartialUser{
		UUID:     newUuid,
		Username: signUpParams.UserName,
	}, nil
}

func (u *userService) GetUser(ctx context.Context, userUUID uuid.UUID) (*userRepo.User, error) {
	return u.userRepo.GetUserByUuid(ctx, userUUID)
}

func (u *userService) GetUserWithOrg(ctx context.Context, userUUID uuid.UUID, orgUUID uuid.UUID) (*userRepo.GetUserWithOrgRow, error) {
	return u.userRepo.GetUserWithOrg(ctx, &userRepo.GetUserWithOrgParams{
		UserUuid: userUUID,
		OrgUuid:  orgUUID,
	})
}

func (u *userService) GetUserEvents(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserEventsRow, error) {
	return u.userRepo.GetUserEvents(ctx, userUUID)
}

func (u *userService) GetUserOrganizations(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserOrganizationsRow, error) {
	return u.userRepo.GetUserOrganizations(ctx, userUUID)
}

func (u *userService) GetEvents(ctx context.Context, searchParams *dto.EventSearch) ([]*userRepo.GetEventsWithOrganizationsRow, error) {
	params := &userRepo.GetEventsWithOrganizationsParams{
		DateFloor:          time.Now(),
		FilterOrganization: false,
		OrgUuid:            [16]byte{},
	}

	if searchParams.DateFloor != nil {
		params.DateFloor = *searchParams.DateFloor
	}

	if searchParams.OrganizationID != nil {
		orgUuid, err := uuid.FromBytes([]byte(*searchParams.OrganizationID))
		if err != nil {
			return nil, err
		}
		params.FilterOrganization = true
		params.OrgUuid = orgUuid
	}

	return u.userRepo.GetEventsWithOrganizations(ctx, params)
}

func (u *userService) CreateEvent(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, eventParams *dto.EventUpsert) error {
	// TODO: validation
	policyNumber := bo.PolicyNumber(orgUser.PoliciesNum)
	if !policyNumber.Can(bo.CREATE_EVENTS) {
		return errors.New("user cannot create events")
	}

	newUuid := uuid.New()

	imageUrl, err := u.storage.Put(ctx, fmt.Sprintf("users/%s.webp", newUuid.String()), &objectstore.Object{
		Data: eventParams.EventImage,
	})
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	if err := u.userRepo.InsertEvent(ctx, &userRepo.InsertEventParams{
		ID:               newUuid,
		EventName:        eventParams.EventName,
		EventDate:        eventParams.EventDate,
		EventLocation:    eventParams.EventLocation,
		EventDescription: eventParams.EventDescription,
		EventImageUrl:    imageUrl,
		// At this point orgfk should be vetted
		OrganizationFk: uuid.Must(uuid.Parse(eventParams.OrganizationID)),
	}); err != nil {
		log.Errorf("%v", err)
		return err
	}

	return nil
}

func (u *userService) CreateInvite(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, inviteParams *dto.InviteCreate) error {
	// TODO: validation
	policyNumber := bo.PolicyNumber(orgUser.PoliciesNum)
	if !policyNumber.Can(bo.CREATE_INVITES) {
		return errors.New("user cannot create invites")
	}

	/* org user is the user creating the invite
	now we have to find the org_user target */

	targetUserUuid, err := uuid.Parse(inviteParams.UserUuid)
	if err != nil {
		return err
	}

	targetEntityUuid, err := uuid.Parse(inviteParams.EntityUuid)
	if err != nil {
		return err
	}

	if err := u.userRepo.InsertInvite(ctx, &userRepo.InsertInviteParams{
		ID:             uuid.New(),
		ExpirationDate: inviteParams.ExpirationDate,
		Capacity:       int32(inviteParams.Capacity),
		UserFk:         targetUserUuid,
		OrgFk:          orgUser.OrganizationFk,
		EntityUuid:     targetEntityUuid,
		EntityType:     0,
	}); err != nil {
		return err
	}

	return nil
}

func (u *userService) GetInvite(ctx context.Context, inviteParams *dto.InviteGet) (*userRepo.Invite, error) {
	uuid, err := uuid.Parse(inviteParams.InviteID)
	if err != nil {
		return nil, err
	}
	return u.userRepo.GetInvite(ctx, uuid)
}

func (u *userService) GetInvitesForOrgUser(ctx context.Context, userUuid uuid.UUID) ([]*userRepo.Invite, error) {
	return u.userRepo.GetInvitesForOrgUser(ctx, userUuid)
}

func (u *userService) GetEvent(ctx context.Context, eventUUID uuid.UUID) (*userRepo.Event, error) {
	return u.userRepo.GetEvent(ctx, eventUUID)
}

func (u *userService) GetOrganization(ctx context.Context, orgUUID uuid.UUID) (*userRepo.Organization, error) {
	return u.userRepo.GetOrganization(ctx, orgUUID)
}

func (u *userService) CreateOrganization(ctx context.Context, orgParams *userRepo.InsertOrganizationParams) error {
	return u.userRepo.InsertOrganization(ctx, orgParams)
}

func (u *userService) CreateOrgUser(ctx context.Context, orgUserParams *userRepo.InsertOrgUserParams) error {
	return u.userRepo.InsertOrgUser(ctx, orgUserParams)
}

// TODO: Remove this and come up with a test only alternative
func (u *userService) DANGER() error {
	return u.userRepo.TruncateAll(context.Background())
}
