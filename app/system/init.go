package system

import (
	"gamma/app/system/auth/ecJwt"

	"github.com/joho/godotenv"
)

var (
	PublicKey  []byte
	PrivateKey []byte
)

func Initialize() {
	err := godotenv.Load()
	if err != nil {
		panic("can't load dot env")
	}
	ecJwt.LoadPrivatePublicKeyPairsDev()
}
