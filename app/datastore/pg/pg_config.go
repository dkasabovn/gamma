package userRepo

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "nhPldb98Rt"
	dbname   = "postgres"
)

var dbCon *sql.DB
var DBSingleton sync.Once

func RwInstance() *sql.DB {
	DBSingleton.Do(func() {
		dbCon = CreateConnection()
	})
	return dbCon
}

func CreateConnection() *sql.DB {
	conString := fmt.Sprintf("port=%d host=%s user=%s password=%s dbname=%s sslmode=disable", port, host, user, password, dbname)
	db, err := sql.Open("postgres", conString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db
}
