package user

import (
	"fmt"
	"net/http"

	"gamma/app/api/core"
	"gamma/app/datastore/users"
	"gamma/app/services/user"
	"gamma/app/system/auth/argon"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	userids struct {
		IDs []primitive.ObjectID `json:"ids"`
	}
)

var (
	coll = user.UserRepo()
)

func getUsers(c echo.Context) error {

	uuids := new(userids)
	if err := c.Bind(uuids); err != nil {
		fmt.Println("Could not prrint")
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}
	
	users, err := coll.GetUsersByUUIDs(c.Request().Context(), uuids.IDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"users" : users,
	}))
}

func getUser(c echo.Context) error {

	gc := core.GetGammaClaims(c)
	uuid := gc.Uuid

	user, err := coll.GetUserByUUID(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}
	return c.JSON(http.StatusOK, core.ApiSuccess(map[string]interface{}{
		"user" : user,
	}))
}

func createUser(c echo.Context) error {

	newUser := new(users.User)
	if err := c.Bind(newUser); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	if err := newUser.ValidNewUser(); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiConverter(map[string]interface{}{
			"error" : err.Error(),
		},0))
	}

	var err error
	newUser.Password, err = argon.PasswordToHash(newUser.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusInternalServerError))
	}

	result, err := coll.CreateUser(c.Request().Context() ,*newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiConverter(map[string]interface{}{
			"error" : err.Error(),
		},0))
	}

	newUser.ID = result.(primitive.ObjectID)
	core.AddTokens(c, *newUser)
	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"uuid" : result,
	}))

}

func login(c echo.Context) error {

	var loginUser users.User

	if err := c.Bind(&loginUser); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	dbUser, err := user.UserRepo().GetUserByEmail(c.Request().Context(), loginUser.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(http.StatusBadRequest))
	}

	match, err := argon.PasswordIsMatch(loginUser.Password, dbUser.Password)
	if !match {
		return c.JSON(http.StatusUnauthorized, core.ApiError(1))
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.ApiError(http.StatusBadRequest))
	}

	core.AddTokens(c, *dbUser)

	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"auth_status": "LOGGED_IN",
	}))
}

// // TODO now needs to update based on email
func updateUser(c echo.Context) error {

	updateData := new(users.User)
	if err := c.Bind(updateData); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(1))
	}
	
	gc := core.GetGammaClaims(c)
	updateData.ID = gc.Uuid

	if err := coll.UpdateUserProfile(c.Request().Context(), *updateData); err != nil {
		return c.JSON(http.StatusBadRequest, core.ApiError(2))
	}
	
	return c.JSON(http.StatusAccepted, core.ApiSuccess(map[string]interface{}{
		"status" : "UPDATED" ,
	}))
}

