package iface

import (
	"context"
	"gamma/app/api/models/dto"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"

	"github.com/google/uuid"
)

type UserService interface {
	SignInUser(ctx context.Context, email, password string) (*bo.PartialUser, error)
	SignUpUser(ctx context.Context, signUpParams *dto.UserSignUp) (*bo.PartialUser, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (*userRepo.User, error)
	GetUserWithOrg(ctx context.Context, userUUID uuid.UUID, orgUUID uuid.UUID) (*userRepo.GetUserWithOrgRow, error)
	GetUserEvents(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserEventsRow, error)
	GetUserOrganizations(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserOrganizationsRow, error)
	GetEvents(ctx context.Context, searchInput *dto.EventSearch) ([]*userRepo.GetEventsWithOrganizationsRow, error)
	CreateEvent(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, eventParams *dto.EventUpsert) error
	CreateOrganization(ctx context.Context, orgParams *userRepo.InsertOrganizationParams) error
	CreateOrgUser(ctx context.Context, orgUserParams *userRepo.InsertOrgUserParams) error
	CreateInvite(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, inviteParams *dto.InviteCreate) error
	GetInvite(ctx context.Context, inviteUUID uuid.UUID) (*userRepo.Invite, error)
	GetEvent(ctx context.Context, eventUUID uuid.UUID) (*userRepo.Event, error)
	GetOrganization(ctx context.Context, orgUUID uuid.UUID) (*userRepo.Organization, error)
	DANGER() error
}
