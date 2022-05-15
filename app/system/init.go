package system

import (
	"gamma/app/system/auth/ecJwt"
)

var (
	PublicKey  []byte
	PrivateKey []byte
)

func Initialize() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Could not find a .env")
	// }
	ecJwt.LoadPrivatePublicKeyPairsDev()
}
