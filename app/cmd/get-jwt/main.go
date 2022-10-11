package main

import (
	"context"
	"gamma/app/api/models/dto"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"gamma/app/system"
	"log"
	"time"

	"github.com/google/uuid"
)

func main() {
	system.Initialize()
	usr, _ := user.GetUserService().SignUpUser(context.Background(), &dto.UserSignUp{
		Email:       "test2@gmail.com",
		PhoneNumber: "5127327045",
		RawPassword: "test",
		FirstName:   "Test",
		LastName:    "Testacular",
		UserName:    "Tester",
	})
	orgUuid := uuid.New()
	user.GetUserService().CreateOrganization(context.Background(), &userRepo.InsertOrganizationParams{
		ID:          orgUuid,
		OrgName:     "Kappa Ligma Balls2",
		City:        "Kappa City",
		OrgImageUrl: "asdfasdf",
	})
	user.GetUserService().CreateOrgUser(context.Background(), &userRepo.InsertOrgUserParams{
		PoliciesNum:    bo.ADMIN,
		UserFk:         usr.UUID,
		OrganizationFk: orgUuid,
	})
	orgUser, err := user.GetUserService().GetUserWithOrg(context.Background(), usr.UUID, orgUuid)
	if err != nil {
		log.Fatal(err)
	}
	var ppl []string
	ppl = append(ppl, usr.UUID.String())
	user.GetUserService().CreateEvent(context.Background(), orgUser, &dto.EventUpsert{
		EventName:           "Woodstock",
		EventDate:           time.Now().Add(100 * time.Hour),
		EventLocation:       "2400 Pearl Street",
		EventDescription:    "Awesome party. Awesome people. Awesome Time.",
		OrganizationID:      orgUuid.String(),
		OrganizationUserIDs: ppl,
	})
}
