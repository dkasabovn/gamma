package main

import (
	"fmt"
	"gamma/app/api/user"

	userDB "gamma/app/datastore/user"

	"github.com/labstack/echo/v4"
)

func main() {

	// privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    // publicKey := &privateKey.PublicKey

    // encPriv, encPub := jwt.EncodeECDSA(privateKey, publicKey)
	// fmt.Println(encPriv)
	// fmt.Println(encPub)

	e := echo.New()

	userDB.MongoUsers()

	e.GET("", func(c echo.Context) error {
		fmt.Println("Printing Cookies")
		fmt.Println(c.Cookies())
		for  _, cookie := range c.Cookies() {
			fmt.Println(cookie.Name)
			fmt.Println(cookie.Value)
	}
	return nil;
	})

	user.JwtRoutes(e)
	user.OpenRoutes(e)
	
	e.Logger.Fatal(e.Start(":8000"))

}
