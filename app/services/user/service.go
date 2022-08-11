package user

import (
	"context"
	"fmt"
	"sync"

	"gamma/app/datastore"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/services/iface"
	"gamma/app/system/auth/argon"
	"gamma/app/system/auth/ecJwt"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

var (
	userOnce            sync.Once
	userServiceInstance iface.UserService
)

type userService struct {
	userRepo *userRepo.Queries
}

func GetUserService() iface.UserService {
	userOnce.Do(func() {
		userServiceInstance = &userService{
			userRepo: userRepo.New(datastore.RwInstance()),
		}
	})
	return userServiceInstance
}

func (u *userService) GetUser(ctx context.Context, uuid string) (*userRepo.User, error) {
	return u.userRepo.GetUserByUuid(ctx, uuid)
}

func (u *userService) GetUserByEmail(ctx context.Context, email string) (*userRepo.User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

func (u *userService) GetOrgUser(ctx context.Context, user_uuid, org_uuid string) (*userRepo.GetOrgUserRow, error) {
	return u.userRepo.GetOrgUser(ctx, &userRepo.GetOrgUserParams{
		UserUuid: user_uuid,
		OrgUuid:  org_uuid,
	})
}

func (u *userService) InsertUser(ctx context.Context, input *userRepo.InsertUserParams) error {
	return u.userRepo.InsertUser(ctx, input)
}

func (u *userService) SignInUser(ctx context.Context, email, password string) (*ecJwt.GammaJwt, error) {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		log.Errorf("Could not get user by Email: %v", email)
		return nil, err
	}

	valid, err := argon.PasswordIsMatch(password, user.PasswordHash)
	if err != nil {
		log.Errorf("error comparing password hashes: %s", err)
		return nil, err
	}

	if valid {
		return ecJwt.GetTokens(ctx, user), nil
	}
	return nil, nil
}

func (u *userService) CreateUser(ctx context.Context, input *userRepo.InsertUserParams) (*ecJwt.GammaJwt, error) {
	hash, err := argon.PasswordToHash(input.PasswordHash)
	if err != nil {
		log.Errorf("could not generate hash: %s", err)
		return nil, err
	}

	input.Uuid = uuid.New().String()
	input.PasswordHash = hash

	if err := u.InsertUser(ctx, input); err != nil {
		// this error should already be logged by InsertUser method
		return nil, err
	}
	return ecJwt.GetTokens(ctx, &userRepo.User{
		Uuid:      input.Uuid,
		ImageUrl:  input.ImageUrl,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}), nil
}

func (u *userService) GetUserOrganizations(ctx context.Context, userId int32) ([]*userRepo.GetUserOrganizationsRow, error) {
	return u.userRepo.GetUserOrganizations(ctx, userId)
}

func (u *userService) GetOrganizationEvents(ctx context.Context, orgUuid string) ([]*userRepo.Event, error) {
	return u.userRepo.GetOrganizationEvents(ctx, orgUuid)
}

func (u *userService) GetUserEvents(ctx context.Context, userId int) ([]*userRepo.GetUserEventsRow, error) {
	return u.userRepo.GetUserEvents(ctx, int32(userId))
}

func (u *userService) GetEvents(ctx context.Context) ([]*userRepo.GetEventsRow, error) {
	return u.userRepo.GetEvents(ctx)
}

func (u *userService) CreateEvent(ctx context.Context, eventParams *userRepo.InsertEventParams) error {
	return u.userRepo.InsertEvent(ctx, eventParams)
}

func (u *userService) CreateOrganization(ctx context.Context, orgParams *userRepo.InsertOrganizationParams) (int32, error) {
	return u.userRepo.InsertOrganization(ctx, orgParams)
}

func (u *userService) CreateOrgUser(ctx context.Context, orgUserParams *userRepo.InsertOrgUserParams) error {
	return u.userRepo.InsertOrgUser(ctx, orgUserParams)
}

func (u *userService) SearchEvents(ctx context.Context, filter string) ([]*userRepo.SearchEventsRow, error) {
	return u.userRepo.SearchEvents(ctx, fmt.Sprintf("'%% %s %%'", filter))
}

func (u *userService) CreateInvite(ctx context.Context, inviteParams *userRepo.InsertInviteParams) error {
	return u.userRepo.InsertInvite(ctx, inviteParams)
}

func (u *userService) GetOrgUserInvites(ctx context.Context, params *userRepo.GetOrgUserInvitesParams) ([]*userRepo.Invite, error) {
	return u.userRepo.GetOrgUserInvites(ctx, params)
}

func (u *userService) GetInvite(ctx context.Context, inviteUuid string) (*userRepo.Invite, error) {
	return u.userRepo.GetInvite(ctx, inviteUuid)
}

func (u *userService) GetEvent(ctx context.Context, eventUuid string) (*userRepo.Event, error) {
	return u.userRepo.GetEvent(ctx, eventUuid)
}

func (u *userService) GetOrganization(ctx context.Context, orgUuid string) (*userRepo.Organization, error) {
	return u.userRepo.GetOrganization(ctx, orgUuid)
}
