package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

	privateKeyName = "PRIVATE_JWT"
	publicKeyName = "PUBLIC_JWT"

	privateKeyPem = "PRIVATE KEY"
	publicKeyPem = "PUBLIC KEY"

)

type UserClaims struct {
	Email   string `bson:"email"`
	UUID    primitive.ObjectID `bson:"_id" json:"uuid"`
	jwt.StandardClaims
}

var (

	KeyFunc = func (*jwt.Token) (interface{}, error) {
		return GetPublicKey(), nil
	}

	// config for JWT tokens specail token name
	CustomJwtConfig = middleware.JWTConfig{
		Claims:      &UserClaims{},
		TokenLookup: fmt.Sprintf("cookie:%s", AccessName),

	}


)

func GetPublicKey() *ecdsa.PublicKey {
	public := user.EnvVariable(publicKeyName)
	private := user.EnvVariable(privateKeyName)

	_, publicKey := DecodeECDSA(private, public)
	return publicKey
}

func getPrivateKey() *ecdsa.PrivateKey {
	public := user.EnvVariable(publicKeyName)
	private := user.EnvVariable(privateKeyName)

	privateKey, _ := DecodeECDSA(private, public)
	return privateKey
}

func GetEmail(cookie http.Cookie) (string, error) {
	token, err := jwt.ParseWithClaims(cookie.Value, &UserClaims{}, KeyFunc)
	if err != nil{
		return "", err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", errors.New("no claims")
	}
	if claims.Email == "" {
		return "", errors.New("no email in Claims")
	}

	return claims.Email , nil
}

func TokenCookie(name, token string, expiration time.Time) *http.Cookie {
	//sets any type of token
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.

	return cookie
}

func generateToken(claim UserClaims, expireTime time.Time) (string, time.Time, error) {
	// generates user token
	claim.StandardClaims = jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		}
	
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
	tokenString, err := token.SignedString(getPrivateKey())

	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expireTime, nil
}

func GenerateAccessToken(claim UserClaims) (string, time.Time, error) {
	expireTime := time.Now().Add(1 * time.Hour)
	return generateToken(claim, expireTime)
}

func GenerateRefreshToken(claim UserClaims) (string, time.Time, error) {
	expireTime := time.Now().Add(72 * time.Hour)
	return generateToken(claim , expireTime)
}

func GenerateServiceToken(claim UserClaims) (string, time.Time, error) {
	expireTime := time.Now().Add(5 * time.Minute)
	return generateToken(claim , expireTime)
}

func EncodeECDSA(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
    x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: privateKeyPem , Bytes: x509Encoded})

    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: publicKeyPem, Bytes: x509EncodedPub})

    return string(pemEncoded), string(pemEncodedPub)
}

func DecodeECDSA(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*ecdsa.PublicKey)

    return privateKey, publicKey
}

func GetClaims(cookie *http.Cookie) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(cookie.Value, &UserClaims{}, KeyFunc);
	if err != nil {
		return nil, errors.New("could not parse token")
	}

	return token.Claims.(*UserClaims), nil
}


