package iface

import (
	"context"
	"gamma/app/domain/bo"
	"gamma/app/system/auth/ecJwt"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, password, email, phone_number, firstName, lastName, userName, image_url string) (*ecJwt.GammaJwt, error)
	SignInUser(ctx context.Context, email, password string) (*ecJwt.GammaJwt, error)
	GetUser(ctx context.Context, uuid string) (*bo.User, error)
	GetUserOrganizations(ctx context.Context, userId int) ([]bo.OrganizationUserJoin, error)
	InsertEventByOrganization(ctx context.Context, orgUuid string, event *bo.Event) (*bo.Event, error)
	CreateInvite(ctx context.Context, orgUser bo.OrgUser, expirationDate time.Time, useLimit int, policy bo.InvitePolicy) (string, error)
	AcceptInvite(ctx context.Context, userEmail string, inviteUuid string) (bool, error)
}
