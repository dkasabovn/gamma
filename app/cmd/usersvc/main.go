package main

import (
	"gamma/app/api/user"
	userDB "gamma/app/datastore/user"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	userDB.MongoUsers()


	user.JwtRoutes(e)
	user.OpenRoutes(e)
	
	e.Logger.Fatal(e.Start(":8000"))

}

// func noAuth(c echo.Context) error {

	

// 	user := userDB.User{
// 		FirstName: "Jazmin",
// 		LastName: "Lowe",
// 		Gender: "F",
// 		Email: "Malinda_Hessel80@yahoo.com",
// 		Bio: "Voluptatibus sapiente velit earum ea asperiores atque consequatur aut. Ullam rem autem adipisci ratione sit maxime. Aut unde magnam voluptas ea. Ea eum voluptas temporibus accusamus voluptatibus libero sunt. Maiores et nobis non modi mollitia vitae.",
// 		HashedPassword: "c50c6ac2-1136-4c66-be68-9ea27dac033a",
// 		ImageLinks: []string{"http://placeimg.com/640/480/abstract"},
// 	}

// 	result, err := userDB.MongoUsers().InsertOne(context.TODO(), user)
// 	if err != nil {
// 		if result != nil {
// 			fmt.Printf("[+] inserted: %s", result.InsertedID)
// 		}
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	return c.JSON(http.StatusAccepted,  result.InsertedID)
// }
