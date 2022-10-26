package ecJwt

import (
	"crypto/ecdsa"
	"sync"
	"time"

	"gamma/app/system"
	"gamma/app/system/log"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type GammaClaims struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	jwt.StandardClaims
}

var (
	keysOnce   sync.Once
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
)

func loadKeys() {
	keysOnce.Do(func() {
		var err error
		privateKey, err = jwt.ParseECPrivateKeyFromPEM([]byte(system.GetConfig().PrivateKey))
		if err != nil {
			panic("can't load private key")
		}
		publicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(system.GetConfig().PublicKey))
		if err != nil {
			panic("can't load public key")
		}
	})
}

func ECDSASign(claims *GammaClaims) (string, string) {
	loadKeys()
	refreshClaims := &GammaClaims{
		UUID: claims.UUID,
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
		log.Errorf("Could not sign jwt with private key: %v", err)
	}

	signedRefreshToken, err := refreshToken.SignedString(privateKey)
	if err != nil {
		log.Errorf("Could not sign jwt with private key: %v", err)
	}

	return signedAccessToken, signedRefreshToken
}

func ECDSAVerify(tokenStr string) (*jwt.Token, bool) {
	loadKeys()
	token, err := jwt.ParseWithClaims(tokenStr, &GammaClaims{}, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	return token, err == nil
}
