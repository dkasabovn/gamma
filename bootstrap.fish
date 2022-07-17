mage run:eventdb
sleep 10
go run ./app/cmd/get-jwt/main.go
go run ./app/cmd/user/main.go