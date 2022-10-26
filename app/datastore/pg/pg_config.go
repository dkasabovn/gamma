package userRepo

import (
	"context"
	"fmt"
	"gamma/app/system"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var dbCon *pgxpool.Pool
var dbSingleton sync.Once

func RwInstance() *pgxpool.Pool {
	dbSingleton.Do(func() {
		dbCon = CreatePool()
	})
	return dbCon
}

func CreateConnection() *pgx.Conn {
	conString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", system.GetConfig().Username, system.GetConfig().Password, system.GetConfig().Hostname, 5432, "postgres")
	db, err := pgx.Connect(context.Background(), conString)
	if err != nil {
		panic(err)
	}

	err = db.Ping(context.Background())

	if err != nil {
		panic(err)
	}

	return db
}

func CreatePool() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.TODO(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", system.GetConfig().Username, system.GetConfig().Password, system.GetConfig().Hostname, 5432, "postgres"))
	if err != nil {
		panic(err)
	}

	err = pool.Ping(context.Background())

	if err != nil {
		panic(err)
	}

	return pool
}
