package dto

import (
	userRepo "gamma/app/datastore/pg"

	"github.com/google/uuid"
)

type OrgUser struct {
	ID           uuid.UUID
	Email        string
	PhoneNumber  string
	FirstName    string
	LastName     string
	Username     string
	ImageUrl     string
	Validated    bool
	PolicyNumber int
}

func ConvertOganizationUser(user *userRepo.GetOrganizationUsersRow) *OrgUser {
	return &OrgUser{
		ID:           user.ID_2,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		ImageUrl:     user.ImageUrl,
		Validated:    user.Validated,
		PolicyNumber: int(user.PoliciesNum),
	}
}

func ConvertOrganizationUsers(users []*userRepo.GetOrganizationUsersRow) []*OrgUser {
	ret := make([]*OrgUser, len(users))
	for i, v := range users {
		ret[i] = ConvertOganizationUser(v)
	}
	return ret
}
