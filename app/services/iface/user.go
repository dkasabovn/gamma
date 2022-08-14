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
	GetOrgUser(ctx context.Context, userUUID uuid.UUID, orgUUID uuid.UUID) (*userRepo.GetOrgUserRow, error)
	GetUserEvents(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserEventsRow, error)
	GetUserOrganizations(ctx context.Context, userUUID uuid.UUID) ([]*userRepo.GetUserOrganizationsRow, error)
	GetEvents(ctx context.Context, searchInput *dto.EventSearch) ([]*userRepo.GetEventsWithOrganizationsRow, error)
	CreateEvent(ctx context.Context, orgUser *userRepo.GetOrgUserRow, eventParams *dto.EventUpsert) error
	CreateOrganization(ctx context.Context, orgParams *userRepo.InsertOrganizationParams) error
	CreateOrgUser(ctx context.Context, orgUserParams *userRepo.InsertOrgUserParams) error
	DANGER() error
}
