package user

import (
	"gamma/app/datastore/mongodb"
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
			db: mongodb.MongoDB().Collection(mongodb.UserCollectionName),
		}
	})
	return userRepo
}
