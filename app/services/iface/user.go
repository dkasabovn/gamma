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
	UpdateUser(ctx context.Context, updateParams *userRepo.User) error
	GetUser(ctx context.Context, userUUID uuid.UUID) (*userRepo.User, error)
	GetUserWithOrg(ctx context.Context, userUUID uuid.UUID, orgUUID uuid.UUID) (*userRepo.GetUserWithOrgRow, error)
	GetUserEvents(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserEventsRow, error)
	GetUserOrganizations(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserOrganizationsRow, error)
	GetEvents(ctx context.Context, searchInput *dto.EventSearch) ([]*userRepo.GetEventsWithOrganizationsRow, error)
	CreateEvent(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, eventParams *dto.EventUpsert) error
	UpdateEvent(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, eventParams *dto.EventUpdate, eventUUID uuid.UUID) error
	DeleteEvent(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, eventUUID uuid.UUID) error
	CreateOrganization(ctx context.Context, orgParams *userRepo.InsertOrganizationParams) error
	CreateOrgUser(ctx context.Context, orgUserParams *userRepo.InsertOrgUserParams) error
	CreateInvite(ctx context.Context, orgUser *userRepo.GetUserWithOrgRow, inviteParams *dto.InviteCreate) error
	GetInvite(ctx context.Context, inviteParams *dto.InviteGet) (*userRepo.Invite, error)
	GetInvitesForOrgUser(ctx context.Context, userUuid uuid.UUID) ([]*userRepo.Invite, error)
	GetEvent(ctx context.Context, eventUUID uuid.UUID) (*userRepo.Event, error)
	GetOrganization(ctx context.Context, orgUUID uuid.UUID) (*userRepo.Organization, error)
	DANGER() error
}
