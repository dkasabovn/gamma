package auth

import (
	"fmt"
	"gamma/app/datastore/user"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// name of tokens in request
	AccessName  = "access-token"
	RefreshName = "refresh-token"

	// keys
	jwtSecretKey  = "SECRET"
	jwtRefreshKey = "REFRESH"
)

type Claims struct {
	ID     primitive.ObjectID `bson:"_id"`
	jwt.StandardClaims
}

var (
	ClaimID = "_id"
	// config for JWT tokens specail token name
	CustomJwtConfig = middleware.JWTConfig{
		Claims:      &Claims{},
		SigningKey:  []byte(jwtSecretKey),
		TokenLookup: fmt.Sprintf("cookie:%s", AccessName),
	}
)

func GetJwtSecret() string {
	return user.EnvVariable("SECRET_JWT")
}

func GetJwtRefresh() string {
	return user.EnvVariable("REFRESH_JWT")
}

func TokenCookie(name, token string, expiration time.Time) *http.Cookie {
	//sets any type of token
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	return cookie
}

func UserCookie(id primitive.ObjectID, expireTime time.Time) *http.Cookie {
	// sets user
	cookie := new(http.Cookie)
	cookie.Name = ClaimID
	cookie.Value = id.Hex()
	cookie.Expires = expireTime
	cookie.Path = "/"

	return cookie
}

func generateToken(id primitive.ObjectID, expireTime time.Time, secret []byte) (string, time.Time, error) {
	// generates user token
	claim := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expireTime, nil
}

func GenerateAccessToken(id primitive.ObjectID) (string, time.Time, error) {
	expireTime := time.Now().Add(1 * time.Hour)
	return generateToken(id, expireTime, []byte(jwtSecretKey))
}

func GenerateRefreshToken(id primitive.ObjectID) (string, time.Time, error) {
	expireTime := time.Now().Add(72 * time.Hour)
	return generateToken(id, expireTime, []byte(jwtRefreshKey))
}
