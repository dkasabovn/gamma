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

func (u *userService) InsertUser(ctx context.Context, uuid string, email string, hash string, firstName string, lastName string, userName string) error {
	return u.userRepo.InsertUser(ctx, uuid, email, hash, firstName, lastName, userName)
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

func (u *userService) CreateUser(ctx context.Context, email, password, firstName, lastName string, userName string) (*ecJwt.GammaJwt, error) {
	hash, err := argon.PasswordToHash(password)
	if err != nil {
		log.Errorf("could not generate hash: %s", err)
		return nil, err
	}
	uuid := uuid.New()
	if err := u.InsertUser(ctx, uuid.String(), email, hash, firstName, lastName, userName); err != nil {
		// this error should already be logged by InsertUser method
		return nil, err
	}
	return ecJwt.GetTokens(ctx, uuid.String(), email, userName, "https://tinyurl.com/monkeygamma"), nil
}

func (u *userService) GetOrgUserEvents(ctx context.Context, user *bo.User) ([]bo.Event, error) {
	if !user.OrgUserFk.Valid {
		return nil, errors.New("org user fk is invalid")
	}
	return u.userRepo.GetOrgUserEvents(ctx, user.OrgUserFk)
}

func (u *userService) GetUserEvents(ctx context.Context, userId int) ([]bo.Event, error) {
	return u.userRepo.GetUserEvents(ctx, userId)
}
