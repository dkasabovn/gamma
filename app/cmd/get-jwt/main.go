package main

import (
	"context"
	"fmt"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/system"
	"gamma/app/system/auth/ecJwt"
)

func main() {
	system.Initialize()
	jwt := ecJwt.GetTokens(context.Background(), &userRepo.User{
		Uuid:      "testing",
		ImageUrl:  "testing",
		FirstName: "testing",
		LastName:  "testing",
	})

	fmt.Printf("%s\n", jwt.BearerToken)
}
