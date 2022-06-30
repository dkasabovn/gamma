package system

import (
	"gamma/app/system/auth/ecJwt"
	"os"

	"github.com/joho/godotenv"
)

var (
	PublicKey   []byte
	PrivateKey  []byte
	ENVIRONMENT string
)

func Initialize() {
	err := godotenv.Load()
	if err != nil {
		panic("can't load dot env")
	}
	ENVIRONMENT = os.Getenv("ENVIRONMENT")
	ecJwt.LoadPrivatePublicKeyPairsDev()
}
