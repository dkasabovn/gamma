package main

import (
	user "gamma/app/api/user"
	"gamma/app/system"
)

func main() {
	system.Initialize()
	user.StartAPI(":6969")
}
