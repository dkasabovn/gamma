package main

import (
	"context"
	"fmt"
	"gamma/app/datastore"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"gamma/app/system"
	"gamma/app/system/auth/ecJwt"
	"time"
)

func main() {
	system.Initialize()

	err := userRepo.New(datastore.RwInstance()).TruncateAll(context.TODO())

	if err != nil {
		fmt.Printf("Error\n")
	}

	tokens, err := user.GetUserService().CreateUser(context.Background(), &userRepo.InsertUserParams{
		Uuid:         "",
		Email:        "bigtest@gmail.com",
		PasswordHash: "testing",
		PhoneNumber:  "6101231234",
		FirstName:    "Big",
		LastName:     "Man",
		ImageUrl:     "https://media.npr.org/assets/img/2017/09/12/macaca_nigra_self-portrait-3e0070aa19a7fe36e802253048411a38f14a79f8-s1100-c50.jpg",
		Validated:    true,
		RefreshToken: "",
	})

	if err != nil {
		panic(err)
	}

	orgId, err := user.GetUserService().CreateOrganization(context.Background(), &userRepo.InsertOrganizationParams{
		OrgName:     "Kappa Ligma",
		City:        "Ligma City",
		Uuid:        "howdypartner",
		OrgImageUrl: "https://media.npr.org/assets/img/2017/09/12/macaca_nigra_self-portrait-3e0070aa19a7fe36e802253048411a38f14a79f8-s1100-c50.jpg",
	})

	if err != nil {
		panic(err)
	}

	tkn, _ := ecJwt.ECDSAVerify(tokens.BearerToken)
	uuid := tkn.Claims.(*ecJwt.GammaClaims).Uuid

	user_obj, err := user.GetUserService().GetUser(context.Background(), uuid)

	if err != nil {
		panic(err)
	}

	err = user.GetUserService().CreateOrgUser(context.Background(), &userRepo.InsertOrgUserParams{
		PoliciesNum:    int32(bo.Create(bo.OWNER)),
		UserFk:         user_obj.ID,
		OrganizationFk: orgId,
	})

	err = user.GetUserService().CreateEvent(context.Background(), &userRepo.InsertEventParams{
		EventName:        "Ligma",
		EventDate:        time.Now(),
		EventLocation:    "Here",
		EventDescription: "howdy partner",
		Uuid:             "poggy woggies",
		EventImageUrl:    "https://media.npr.org/assets/img/2017/09/12/macaca_nigra_self-portrait-3e0070aa19a7fe36e802253048411a38f14a79f8-s1100-c50.jpg",
		OrganizationFk:   orgId,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Access: %s\n\n\nRefresh: %s", tokens.BearerToken, tokens.RefreshToken)
}
