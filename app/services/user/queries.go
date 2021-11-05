package user

import (
	"context"
	"gamma/app/datastore/mongodb/models"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *UserRepository) GetUserByUUID(ctx context.Context, uuid primitive.ObjectID) (*models.User, error) {
	return models.UserByUUID(ctx, u.db, uuid)
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return models.UserByEmail(ctx, u.db, email)
}

func (u *UserRepository) GetUsersByUUIDs(ctx context.Context, uuids []primitive.ObjectID) ([]models.User, error) {
	return models.UsersByUUIDs(ctx, u.db, uuids)
}

func (u *UserRepository) CreateUser(ctx context.Context, user models.User) (interface{}, error) {
	res, err := models.NewUser(ctx, u.db, user)
	if err != nil {
		log.Printf("%+v", err)
		return primitive.NilObjectID, err
	}
	return res.InsertedID, err
}

func (u *UserRepository) UpdateUserProfile(ctx context.Context, user models.User) error {
	return models.UpdateUser(ctx, u.db, user)
}
