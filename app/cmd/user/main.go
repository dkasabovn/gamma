package main

import (
	user "gamma/app/api/user"
	"gamma/app/system"
)

func main() {
	system.Initialize()
	api := user.API()
	api.Start(":6969")
}
