package ecJwt

import (
	"context"
	"crypto/ecdsa"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type GammaClaims struct {
	Uuid string `json:"uuid"`
	jwt.StandardClaims
}

type GammaJwt struct {
	BearerToken  string `json:"bearer_token" header:"Authorization"`
	RefreshToken string `json:"refresh_token" param:"refresh_token"`
}

// TODO !!! Handle reading private, public key files
func ECDSASign(claims *GammaClaims) (string, string) {
	privateKeyString, _ := os.ReadFile("private-key.pem")

	var privateKey *ecdsa.PrivateKey
	var err error
	if privateKey, err = jwt.ParseECPrivateKeyFromPEM([]byte(privateKeyString)); err != nil {
		log.Fatalf("Unable to parse ECDSA private key: %v", err)
	}

	refreshClaims := &GammaClaims{
		Uuid: claims.Uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    "auth.gamma",
			Audience:  "user.gamma",
		},
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		refreshClaims,
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
	publicKeyString, _ := os.ReadFile("public-key.pem")
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

func GetTokens(ctx context.Context, userUuid string) *GammaJwt {
	claims := &GammaClaims{
		Uuid: userUuid,
	}

	accessToken, refreshToken := ECDSASign(claims)
	return &GammaJwt{
		BearerToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
