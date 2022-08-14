package dto

import (
	userRepo "gamma/app/datastore/pg"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	PhoneNumber string
	FirstName   string
	LastName    string
	Username    string
	ImageUrl    string
	Validated   bool
}

func ConvertUser(user *userRepo.User) *User {
	return &User{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		ImageUrl:    user.ImageUrl,
		Validated:   user.Validated,
	}
}
