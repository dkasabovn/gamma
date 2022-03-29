package system

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Initialize() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Could not find a .env")
	// }
	path, _ := os.Getwd()
	fmt.Printf("DIR: %s", path)
	loadPrivatePublicKeyPairs()
}
func loadPrivatePublicKeyPairs() {
	private_key, err := ioutil.ReadFile("private-key.pem")
	if err != nil {
		panic("Could not open private-key file")
	}
	os.Setenv("PRIVATE_KEY", string(private_key))
	public_key, err := ioutil.ReadFile("public-key.pem")
	if err != nil {
		panic("Could not open public-key file")
	}
	os.Setenv("PUBLIC_KEY", string(public_key))
}
