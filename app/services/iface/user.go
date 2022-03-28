package iface

import (
	"context"
	"gamma/app/system/auth/ecJwt"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName, userName string) (*ecJwt.GammaJwt, error)
	SignInUser(ctx context.Context, email, password string) (*ecJwt.GammaJwt, error)
}
