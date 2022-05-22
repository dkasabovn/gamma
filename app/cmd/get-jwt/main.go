package main

import (
	"context"
	"fmt"
	"gamma/app/datastore"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/services/user"
	"gamma/app/system"
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
		fmt.Printf("Error\n")
	}

	fmt.Printf("Access: %s\n\n\nRefresh: %s", tokens.BearerToken, tokens.RefreshToken)
}
