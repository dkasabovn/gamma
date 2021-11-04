package user

import (
	"context"
	"gamma/app/datastore/users"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *UserRepository) GetUserByUUID(ctx context.Context, uuid primitive.ObjectID) (*users.User, error) {
	return users.UserByUUID(ctx, u.db, uuid)
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	return users.UserByEmail(ctx, u.db, email)
}

func (u *UserRepository) GetUsersByUUIDs(ctx context.Context, uuids []primitive.ObjectID) ([]users.User, error) {
	return users.UsersByUUIDs(ctx, u.db, uuids)
}

func (u *UserRepository) CreateUser(ctx context.Context, user users.User) (interface{}, error) {
	res, err := users.NewUser(ctx, u.db, user)
	if err != nil {
		log.Printf("%+v", err)
		return primitive.NilObjectID, err
	}
	return res.InsertedID, err
}

func (u *UserRepository) UpdateUserProfile(ctx context.Context, user users.User) error {
	return users.UpdateUser(ctx, u.db, user)
}
