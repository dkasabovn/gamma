package ecJwt

import (
	"crypto/ecdsa"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type GammaClaims struct {
	Uuid  string `json:"uuid"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func ECDSASign(claims *GammaClaims) (string, string) {
	privateKeyString, _ := os.ReadFile("private-key.pem")

	var privateKey *ecdsa.PrivateKey
	var err error
	if privateKey, err = jwt.ParseECPrivateKeyFromPEM([]byte(privateKeyString)); err != nil {
		log.Fatalf("Unable to parse ECDSA private key: %v", err)
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    "auth.gamma",
			Audience:  "user.gamma",
		},
	)

	claims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
	claims.Issuer = "auth.gamma"
	claims.Audience = "user.gamma"

	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		claims,
	)

	signedAccessToken, err := token.SignedString(privateKey)

	if err != nil {
		log.Fatalf("Could not sign jwt with private key: %v", err)
	}

	signedRefreshToken, err := refreshToken.SignedString(privateKey)

	if err != nil {
		log.Fatalf("Could not sign jwt with private key: %v", err)
	}

	return signedAccessToken, signedRefreshToken
}

func ECDSAVerify(tokenStr string) (*jwt.Token, bool) {
	publicKeyString := os.Getenv("PUBLIC_KEY")
	var publicKey *ecdsa.PublicKey
	var err error
	if publicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(publicKeyString)); err != nil {
		log.Fatalf("Unable to parse ECDSA public key: %v", err)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &GammaClaims{}, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	return token, err == nil
}
