package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	refresh "gamma/app/api/core"
	userDB "gamma/app/datastore/user"
	jwt "gamma/app/system/auth"
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

var (
	displayData = bson.M{"hashedPassword": 0, "device": 0}
)



func validateNewUser(user *userDB.User) error {
	if user.HashedPassword == "" {
		fmt.Println("password")
		return errors.New("missing password")
	}
	if user.Email == "" {
		fmt.Println("email")

		return errors.New("missing email")
	} 
	if user.FirstName == "" {
		fmt.Println("firstname")

		return errors.New("missing firstname")
	}
	if user.DisplayName == "" {
		fmt.Println("display")

		return errors.New("missing DisplayName")
	}
	if user.Bio == "" {
		fmt.Println("bio")

		return errors.New("missing bio")
	}

	return nil

}

func getUsers(c echo.Context) error {
	ids := new(userids)
	
	if err := c.Bind(ids); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	filter := bson.M{"_id": bson.M{"$in": ids.IDs}}
	opts := options.Find().SetProjection(displayData)
	
	cursor, err := userDB.MongoUsers().Find(c.Request().Context(), filter, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var users []userDB.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {
	
	email, err := getEmail(c)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err.Error())
	}

	filter := bson.M{"email": email}
	opts := options.FindOne().SetProjection(displayData)
	
	var dbUser userDB.User
	if err := userDB.MongoUsers().FindOne(c.Request().Context(), filter, opts).Decode(&dbUser); err != nil {
		return c.JSON(http.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(http.StatusAccepted, dbUser)
}

func createUser(c echo.Context) error {

	user := new(userDB.User)
	
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "binding issue")
	}

	if err := validateNewUser(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var err error
	user.HashedPassword, err = argon.PasswordToHash(user.HashedPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "HASHING ISSUE")
	}

	result, err := userDB.MongoUsers().InsertOne(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	fmt.Printf("[+] inserted: %s", result.InsertedID)
	return c.JSON(http.StatusAccepted,  result.InsertedID)
}

func login(c echo.Context) error {

	var user userDB.User
	
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}


	var dbUser userDB.User
	filter := bson.M{"email": user.Email}

	if err := userDB.MongoUsers().FindOne(c.Request().Context(), filter).Decode(&dbUser); err != nil {
		return c.JSON(http.StatusServiceUnavailable, err)
	}

	

	match, err := argon.PasswordIsMatch(user.HashedPassword, dbUser.HashedPassword)
	if !match {
		return c.JSON(http.StatusUnauthorized, "Incorrect Password")
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "HASHING error")
	}
	
	claims := &jwt.UserClaims{
		Email:  dbUser.Email,
		UUID: dbUser.ID,
	}

	if err := refresh.UpdateTokens(*claims, c); err != nil {
		return c.JSON(http.StatusUnauthorized, "Token is invalid")
	}

	return c.JSON(http.StatusAccepted, "Logged in")
}

// TODO now needs to update based on email
func updateUser(c echo.Context) error {

	userCookie, _ := c.Cookie(jwt.RefreshName)
	email, _ := jwt.GetEmail(*userCookie)

	updateData := new(userDB.User)
	if err := c.Bind(updateData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	updates := bson.M{
		"$set":  bson.M{
        	"imageLinks":  updateData.ImageLinks,
        	"bio": updateData.Bio,
    }}

	result, err := userDB.MongoUsers().UpdateByID(c.Request().Context(), email, updates)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err.Error())
	}
	if result.ModifiedCount == 0 {
		return c.JSON(http.StatusBadRequest, "User does not exist")
	}

	return c.JSON(http.StatusAccepted, "updated")
}

func getEmail(c echo.Context) (string, error) {
	accessCookie, err := c.Cookie(jwt.AccessName)
	if err != nil {
		return "", err
	}

	claims, err := jwt.GetClaims(accessCookie)
	if err != nil {
		return "", err
	}

	return claims.Email, nil


	
}
