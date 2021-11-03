package user

import (
	"gamma/app/datastore/users"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userCollection *mongo.Collection
	
	userSingle sync.Once
	userRepo *UserRepository


)

type UserRepository struct {
	db *mongo.Collection
}

func UserRepo() *UserRepository {
	userSingle.Do(func() {
		userRepo = &UserRepository{
			db: users.MongoDB().Collection(users.UserCollectionName),
		}
	})
	return userRepo
}
