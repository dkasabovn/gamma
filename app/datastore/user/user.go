package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	User struct {
		ID 			primitive.ObjectID  `bson:"_id,omitempty"`
		DisplayName     string              `bson:"displayName,omitempty"`
		FirstName 		string 			    `bson:"firstName,omitempty"`
		LastName 		string 			    `bson:"lastName,omitempty"`
		Gender 			string               `bson:"gender,omitempty"`	
		Email 			string              `bson:"email,omitempty"`
		Bio 			string              `bson:"bio,omitempty"`
		HashedPassword 	string              `bson:"hashedPassword,omitempty"`
		ImageLinks 		[]string            `bson:"imageLinks,omitempty"`
	}
)


