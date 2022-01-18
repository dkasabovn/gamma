package user

import (
	"context"
	"errors"
	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/domain/definition"
	"gamma/app/services/iface"
	"sync"

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

func (u *userService) InsertUser(ctx context.Context, uuid string, email string, firstName string, lastName string) error {
	return u.userRepo.InsertUser(ctx, uuid, email, firstName, lastName)
}

func (u *userService) GetOrgUserEvents(ctx context.Context, user *bo.User) ([]bo.Event, error) {
	if !user.OrgUserFk.Valid {
		return nil, errors.New("org user fk is invalid")
	}
	return u.userRepo.GetOrgUserEvents(ctx, int(user.OrgUserFk.Int64))
}

func (u *userService) GenerateInvite(ctx context.Context, eventUuid string) (string, error) {
	inviteId := uuid.New()
	err := u.userRepo.InsertInvite(ctx, inviteId.String(), eventUuid)
	if err != nil {
		log.Errorf("could not insert invite: %s", err.Error())
		return "", err
	}
	return inviteId.String(), nil
}

func (u *userService) AcceptInvite(ctx context.Context, user *bo.User, inviteUuid string) error {
	invite, err := u.userRepo.GetInvite(ctx, inviteUuid)
	if err != nil {
		log.Errorf("could not get invite: %s", err.Error())
		return err
	}

	event, err := u.userRepo.GetEvent(ctx, invite.EventUuid)
	if err != nil {
		log.Errorf("could not get event by invite eventuuid: %s", err.Error())
		return err
	}

	err = u.userRepo.InsertUserEvent(ctx, user.Id, event.Id)
	if err != nil {
		log.Errorf("could not add event to user: %s", err)
		return err
	}

	return nil
}

func (u *userService) GetUserEvents(ctx context.Context, userId int) ([]bo.Event, error) {
	return u.userRepo.GetUserEvents(ctx, userId)
}
