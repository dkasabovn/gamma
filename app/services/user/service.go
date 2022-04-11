package user

import (
	"context"
	"errors"
	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/domain/definition"
	"gamma/app/services/iface"
	"gamma/app/system/auth/argon"
	"gamma/app/system/auth/ecJwt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

var (
	userOnce            sync.Once
	userServiceInstance iface.UserService
)

type userService struct {
	userRepo definition.UserRepository
}

func GetUserService() iface.UserService {
	userOnce.Do(func() {
		userServiceInstance = &userService{
			userRepo: pg.GetUserRepo(),
		}
	})
	return userServiceInstance
}

func (u *userService) GetUser(ctx context.Context, uuid string) (*bo.User, error) {
	return u.userRepo.GetUser(ctx, uuid)
}

func (u *userService) GetUserByEmail(ctx context.Context, email string) (*bo.User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

func (u *userService) InsertUser(ctx context.Context, uuid, email, phone_number, hash, firstName, lastName, userName, image_url string) error {
	return u.userRepo.InsertUser(ctx, uuid, email, phone_number, hash, firstName, lastName, userName, image_url)
}

func (u *userService) SignInUser(ctx context.Context, email, password string) (*ecJwt.GammaJwt, error) {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	valid, err := argon.PasswordIsMatch(password, user.PasswordHash)

	if err != nil {
		log.Errorf("error comparing password hashes: %s", err)
		return nil, err
	}

	if valid {
		return ecJwt.GetTokens(ctx, user.Uuid, user.Email, user.UserName, "https://tinyurl.com/monkeygamma"), nil
	}
	return nil, nil
}

func (u *userService) CreateUser(ctx context.Context, password, email, phone_number, firstName, lastName, userName, image_url string) (*ecJwt.GammaJwt, error) {
	hash, err := argon.PasswordToHash(password)
	if err != nil {
		log.Errorf("could not generate hash: %s", err)
		return nil, err
	}
	uuid := uuid.New()
	if err := u.InsertUser(ctx, uuid.String(), email, phone_number, hash, firstName, lastName, userName, image_url); err != nil {
		// this error should already be logged by InsertUser method
		return nil, err
	}
	return ecJwt.GetTokens(ctx, uuid.String(), email, userName, "https://tinyurl.com/monkeygamma"), nil
}

func (u *userService) GetUserOrganizations(ctx context.Context, userId int) ([]bo.OrganizationUserJoin, error) {
	return u.userRepo.GetUserOrganizations(ctx, userId)
}

func (u *userService) InsertEventByOrganization(ctx context.Context, orgUuid string, event *bo.Event) (*bo.Event, error) {
	return u.userRepo.InsertEventByOrganization(ctx, orgUuid, event)
}

func (u *userService) GetUserEvents(ctx context.Context, userId int) ([]bo.Event, error) {
	return u.userRepo.GetUserEvents(ctx, userId)
}

func (u *userService) CreateInvite(ctx context.Context, orgUser bo.OrgUser, expirationDate time.Time, useLimit int, policy bo.InvitePolicy) (string, error)  {
	if !orgUser.CanCreateEvent() {
		return  "", errors.New("user does not have perms")
	}

	inviteUuid := uuid.NewString()
	return inviteUuid, u.userRepo.InsertInvite(ctx, expirationDate, useLimit, inviteUuid, policy)
}

func(u *userService) AcceptInvite(ctx context.Context, userEmail string, inviteUuid string) (bool, error) {
	invite, err := u.userRepo.GetInvite(ctx, inviteUuid)
	if err != nil {
		return false, err
	}

	user, err := u.userRepo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return false, err
	}

	userOrgs, err := u.userRepo.GetUserOrganizations(ctx, user.Id)
	if err != nil {
		return false, err
	}

	if !invite.UserCanAttend(user, userOrgs) {
		return false, nil
	}

	switch invite.Policy.InviteType{
	case bo.InviteToEvent:	
		err = u.userRepo.InsertUserEvent(ctx, user.Id, invite.Policy.InviteTo)
	case bo.InviteToOrg:
		err = u.userRepo.InsertOrgUser(ctx, user.Id, invite.Policy.InviteTo, invite.Policy.PoliciesNum)
	}

	if err != nil {
		return false, err
	}

	return true, u.userRepo.DecrementInvite(ctx, invite.Id)
}





