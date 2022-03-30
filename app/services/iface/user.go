package iface

import (
	"context"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/ecJwt"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName, userName string) (*ecJwt.GammaJwt, error)
	SignInUser(ctx context.Context, email, password string) (*ecJwt.GammaJwt, error)
	GetUser(ctx context.Context, uuid string) (*bo.User, error)
	GetUserOrganizations(ctx context.Context, userId int) ([]bo.OrganizationUser, error)
	InsertEventByOrganization(ctx context.Context, orgUuid string, event *bo.Event) (*bo.Event, error)
}
