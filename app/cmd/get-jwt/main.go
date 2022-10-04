package main

import (
	"context"
	"gamma/app/api/models/dto"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"gamma/app/system"

	"github.com/google/uuid"
)

func main() {
	system.Initialize()
	usr, _ := user.GetUserService().SignUpUser(context.Background(), &dto.UserSignUp{
		Email:       "test@gmail.com",
		PhoneNumber: "5126327045",
		RawPassword: "test",
		FirstName:   "Test",
		LastName:    "Testacular",
		UserName:    "Tester",
	})
	orgUuid := uuid.New()
	user.GetUserService().CreateOrganization(context.Background(), &userRepo.InsertOrganizationParams{
		ID:          orgUuid,
		OrgName:     "Kappa Ligma Balls",
		City:        "Kappa City",
		OrgImageUrl: "asdfasdf",
	})
	user.GetUserService().CreateOrgUser(context.Background(), &userRepo.InsertOrgUserParams{
		PoliciesNum:    bo.ADMIN,
		UserFk:         usr.UUID,
		OrganizationFk: orgUuid,
	})
}
