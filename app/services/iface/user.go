package iface

import (
	"context"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/system/auth/ecJwt"
)

type UserService interface {
	CreateUser(ctx context.Context, input *userRepo.InsertUserParams) (*ecJwt.GammaJwt, error)
	SignInUser(ctx context.Context, email, password string) (*ecJwt.GammaJwt, error)
	GetUser(ctx context.Context, uuid string) (*userRepo.User, error)
	GetOrgUser(ctx context.Context, user_uuid, org_uuid string) (*userRepo.GetOrgUserRow, error)
	GetUserOrganizations(ctx context.Context, userId int32) ([]*userRepo.GetUserOrganizationsRow, error)
	GetEvents(ctx context.Context, userId int) ([]*userRepo.GetEventsRow, error)
	CreateEvent(ctx context.Context, eventParams *userRepo.InsertEventParams) error
	CreateOrganization(ctx context.Context, orgParams *userRepo.InsertOrganizationParams) (int32, error)
	CreateOrgUser(ctx context.Context, orgUserParams *userRepo.InsertOrgUserParams) error
	GetOrganizationEvents(ctx context.Context, orgUuid string) ([]*userRepo.Event, error)
	GetUserEvents(ctx context.Context, userId int) ([]*userRepo.GetUserEventsRow, error)
	CreateInvite(ctx context.Context, inviteParams *userRepo.InsertInviteParams) error
}
