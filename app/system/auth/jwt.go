package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
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
	Email string `json:"email"`
	jwt.StandardClaims
}

var (

	// config for JWT tokens specail token name
	CustomJwtConfig = middleware.JWTConfig{
		Claims:      &Claims{},
		SigningKey:  []byte(jwtSecretKey),
		TokenLookup: fmt.Sprintf("cookie:%s", AccessName),
	}
)

func GetJwtSecret() string {
	return jwtSecretKey
}

func GetJwtRefresh() string {
	return jwtRefreshKey
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

func UserCookie(email string, expireTime time.Time) *http.Cookie {
	// sets user
	cookie := new(http.Cookie)
	cookie.Name = "email"
	cookie.Value = email
	cookie.Expires = expireTime
	cookie.Path = "/"

	return cookie
}

func generateToken(email string, expireTime time.Time, secret []byte) (string, time.Time, error) {
	// generates user token
	claim := &Claims{
		Email: email,
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

func GenerateAccessToken(email string) (string, time.Time, error) {
	expireTime := time.Now().Add(1 * time.Hour)
	return generateToken(email, expireTime, []byte(jwtSecretKey))
}

func GenerateRefreshToken(email string) (string, time.Time, error) {
	expireTime := time.Now().Add(72 * time.Hour)
	return generateToken(email, expireTime, []byte(jwtRefreshKey))
}
