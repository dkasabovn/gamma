package ecJwt

import (
	"crypto/ecdsa"
	"io/ioutil"
	"time"

	"gamma/app/system/log"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	privateKey []byte
	publicKey  []byte
)

type GammaClaims struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	jwt.StandardClaims
}

func LoadPrivatePublicKeyPairsDev() {
	private_key, err := ioutil.ReadFile("private-key.pem")
	if err != nil {
		panic("Could not open private-key file")
	}
	privateKey = private_key
	public_key, err := ioutil.ReadFile("public-key.pem")
	if err != nil {
		panic("Could not open public-key file")
	}
	publicKey = public_key
}

func ECDSASign(claims *GammaClaims) (string, string) {
	var privateKeyS *ecdsa.PrivateKey
	var err error

	if privateKeyS, err = jwt.ParseECPrivateKeyFromPEM(privateKey); err != nil {
		log.Errorf("Unable to parse ECDSA private key: %v", err)
	}

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

	signedAccessToken, err := token.SignedString(privateKeyS)
	if err != nil {
		log.Errorf("Could not sign jwt with private key: %v", err)
	}

	signedRefreshToken, err := refreshToken.SignedString(privateKeyS)
	if err != nil {
		log.Errorf("Could not sign jwt with private key: %v", err)
	}

	return signedAccessToken, signedRefreshToken
}

func ECDSAVerify(tokenStr string) (*jwt.Token, bool) {
	var err error
	var publicKeyS *ecdsa.PublicKey

	if publicKeyS, err = jwt.ParseECPublicKeyFromPEM(publicKey); err != nil {
		log.Errorf("Unable to parse ECDSA public key: %v", err)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &GammaClaims{}, func(t *jwt.Token) (interface{}, error) {
		return publicKeyS, nil
	})
	return token, err == nil
}
