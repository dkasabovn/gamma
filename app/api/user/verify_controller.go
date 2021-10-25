package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	userDB "gamma/app/datastore/user"
	"gamma/app/system/auth/argon"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type(
	userids struct {
		IDs []primitive.ObjectID `bson:"ids"`
	}
)

func validateNewUser(user *userDB.User) error {
	if user.HashedPassword == "" {
		return errors.New("missing password")
	}
	if user.Email == "" {
		return errors.New("missing email")
	} 
	if user.FirstName == "" {
		return errors.New("missing firstname")
	}
	if user.DisplayName == "" {
		return errors.New("missing DisplayName")
	}
	if user.Bio == "" {
		return errors.New("missing bio")
	}

	return nil

}

func getUsers(c echo.Context) error {
	ids := new(userids)
	
	if err := c.Bind(ids); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	filter := bson.M{"_id": bson.M{"$in": ids.IDs}}
	opts := options.Find().SetProjection(bson.M{"hashedPassword": 0})
	
	// var users []userDB.User
	cursor, err := userDB.MongoUsers().Find(context.TODO(), filter, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var users []userDB.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, users)
}

func createUser(c echo.Context) error {

	user := new(userDB.User)
	
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := validateNewUser(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var err error
	user.HashedPassword, err = argon.PasswordToHash(user.HashedPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "HASHING ISSUE")
	}

	result, err := userDB.MongoUsers().InsertOne(context.TODO(), user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	
	fmt.Printf("[+] inserted: %s", result.InsertedID)
	return c.JSON(http.StatusAccepted,  result.InsertedID)
}

func login(c echo.Context) error {

	var user userDB.User
	
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := validateNewUser(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}


	var dbUser userDB.User
	opts := options.FindOne().SetProjection(bson.M{"hashedPassword": 1})
	filter := bson.M{"email": user.Email}

	if err := userDB.MongoUsers().FindOne(context.TODO(), filter, opts).Decode(&dbUser); err != nil {
		fmt.Println("user not found")
		return c.JSON(http.StatusServiceUnavailable, err)
	}
	

	match, err := argon.PasswordIsMatch(user.HashedPassword, dbUser.HashedPassword)
	if !match {
		fmt.Println("poops")
		return c.JSON(http.StatusUnauthorized, "Incorrect Password")
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "HASHING error")
	}
	fmt.Println("Found could match")

	if err := updateTokens(dbUser.Email, c); err != nil {
		return c.JSON(http.StatusUnauthorized, "Token is invalid")
	}

	return c.JSON(http.StatusAccepted, "Logged In")
	

}

