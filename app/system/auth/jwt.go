package auth

import (
	"time"

	"net/http"
	_ "net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	AccessName  = "access-token"
	RefreshName = "refresh-token"

	jwtSecretKey  = "SECRET"
	jwtRefreshKey = "REFRESH"
)

type User struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
}
type Claims struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GetJwtSecret() string {
	return jwtSecretKey
}

func GetJwtRefresh() string {
	return jwtRefreshKey
}

func generateToken(user *User, expireTime time.Time, secret []byte) (string, time.Time, error) {
	// generates user token
	claim := &Claims{
		Email: user.Email,
		UUID:  user.UUID,
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

func generateAccessToken(user *User) (string, time.Time, error) {
	expireTime := time.Now().Add(1 * time.Hour)
	return generateToken(user, expireTime, []byte(jwtSecretKey))
}

func generateRefreshToken(user *User) (string, time.Time, error) {
	expireTime := time.Now().Add(24 * time.Hour)
	return generateToken(user, expireTime, []byte(jwtRefreshKey))
}

func setTokenCookie(name, token string, expiration time.Time, ctx echo.Context) {
	//sets any type of token
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	ctx.SetCookie(cookie)
}

func setUserCookie(user *User, expireTime time.Time, c echo.Context) {
	// sets user
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.UUID
	cookie.Expires = expireTime
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func GenerateTokenAndSetCookies(user *User, ctx echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(AccessName, accessToken, exp, ctx)
	setUserCookie(user, exp, ctx)

	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return nil
	}

	setTokenCookie(RefreshName, refreshToken, exp, ctx)
	return nil
}

func ErrorHandler(err error, ctx echo.Context) error {
	return ctx.Redirect(http.StatusMovedPermanently, "")
}

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Get("user") == nil {
			return next(ctx)
		}

		u := ctx.Get("user").(*jwt.Token)
		claims := u.Claims.(*Claims)

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {

			refreshCookie, err := ctx.Cookie(RefreshName)
			if err == nil && refreshCookie != nil {

				refreshTkn, err := jwt.ParseWithClaims(refreshCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetJwtRefresh()), nil
				})

				if err != nil {
					ctx.Response().Writer.WriteHeader(http.StatusUnauthorized)
				}

				if refreshTkn != nil && refreshTkn.Valid {
					_ = GenerateTokenAndSetCookies(&User{
						UUID: claims.UUID,
					}, ctx)
				}

			}

		}

		return next(ctx)
	}
}
