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
	Email string `json:"email,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Image string `json:"image,omitempty"`
	jwt.StandardClaims
}

type GammaJwt struct {
	BearerToken  string `json:"bearer_token" header:"Authorization"`
	RefreshToken string `json:"refresh_token" param:"refresh_token"`
}

// TODO !!! Handle reading private, public key files
func ECDSASign(claims *GammaClaims) (string, string) {
	var privateKeyString []byte
	var err error
	privateKeyString = []byte(os.Getenv("PRIVATE_KEY"))

	var privateKey *ecdsa.PrivateKey
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
	var publicKeyString []byte
	var err error
	publicKeyString = []byte(os.Getenv("PUBLIC_KEY"))
	
	var publicKey *ecdsa.PublicKey
	if publicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(publicKeyString)); err != nil {
		log.Fatalf("Unable to parse ECDSA public key: %v", err)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &GammaClaims{}, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	return token, err == nil
}

func GetTokens(ctx context.Context, userUuid, userEmail, userName, userImage string) *GammaJwt {
	claims := &GammaClaims{
		Uuid: userUuid,
		Email: userEmail,
		UserName: userName,
		Image: userImage,
	}

	accessToken, refreshToken := ECDSASign(claims)
	return &GammaJwt{
		BearerToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
