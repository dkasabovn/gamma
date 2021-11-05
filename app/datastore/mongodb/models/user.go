package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	standardFilter = bson.M{"hashedPassword": 0, "device": 0}
	findOneOpts    = options.FindOne().SetProjection(standardFilter)
	findOpts       = options.Find().SetProjection(standardFilter)
)

const ()

type (
	User struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"uuid,omitempty"`
		DisplayName string             `bson:"displayName,omitempty" json:"username,omitempty"`
		FirstName   string             `bson:"firstName,omitempty" json:"firstName,omitempty"`
		LastName    string             `bson:"lastName,omitempty" json:"lastName,omitempty"`
		Gender      string             `bson:"gender,omitempty" json:"gender,omitempty"`
		Email       string             `bson:"email,omitempty" json:"email,omitempty"`
		Bio         string             `bson:"bio,omitempty" json:"bio,omitempty"`
		Password    string             `bson:"hashedPassword,omitempty" json:"password,omitempty"`
		ImageLinks  []string           `bson:"imageLinks,omitempty" json:"links,omitempty"`
		Device      string             `bson:"device,omitempty" json:"device,omitempty"`
		CreatedAt   time.Time 				`bson:"createdAt, omitempty" json:"createdAt"`
	}
)

func (u *User) ValidNewUser() error {
	if u.Password == "" {
		fmt.Println("password")
		return errors.New("missing password")
	}
	if u.Email == "" {
		fmt.Println("email")

		return errors.New("missing email")
	}
	if u.FirstName == "" {
		fmt.Println("firstname")

		return errors.New("missing firstname")
	}
	if u.DisplayName == "" {
		fmt.Println("display")

		return errors.New("missing DisplayName")
	}
	if u.Bio == "" {
		fmt.Println("bio")

		return errors.New("missing bio")
	}
	if len(u.ImageLinks) == 0 {
		return errors.New("missing images")
	}
	if u.Device == "" {
		return errors.New("missing device")
	}

	if u.Gender != "" {
		if u.Gender != "M" && u.Gender != "F" && u.Gender != "O" &&  u.Gender != "N" {
			return errors.New("Not valid gender has to be M,F,O,N")
		}
	}

	return nil

}

func UserByUUID(ctx context.Context, coll *mongo.Collection, uuid primitive.ObjectID) (*User, error) {
	user := new(User)
	filter := bson.M{"_id": uuid}

	err := coll.FindOne(ctx, filter, findOneOpts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UserByEmail(ctx context.Context, coll *mongo.Collection, email string) (*User, error) {
	user := new(User)
	filter := bson.M{"email": email}

	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UsersByUUIDs(ctx context.Context, coll *mongo.Collection, ids []primitive.ObjectID) ([]User, error) {
	var users []User

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func NewUser(ctx context.Context, coll *mongo.Collection, user User) (*mongo.InsertOneResult, error) {
	user.CreatedAt = time.Now()
	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateUser(ctx context.Context, coll *mongo.Collection, user User) error {
	updates := bson.M{
		"$set": bson.M{
			"imageLinks": user.ImageLinks,
			"bio":        user.Bio,
		}}
	result, err := coll.UpdateByID(ctx, user.ID, updates)
	if result.ModifiedCount == 0 {
		return errors.New("nothing user found")
	}
	if err != nil {
		return err
	}

	return nil

}
