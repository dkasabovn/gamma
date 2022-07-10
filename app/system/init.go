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
	BUCKET_NAME string
)

func Initialize() {
	err := godotenv.Load()
	if err != nil {
		panic("can't load dot env")
	}
	ENVIRONMENT = os.Getenv("ENVIRONMENT")
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
	ecJwt.LoadPrivatePublicKeyPairsDev()
}
