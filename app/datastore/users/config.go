package users

import (
	"context"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	
	db = "gamma"
	UserCollectionName = "users"


)

var (
	mongoClient *mongo.Client
	mongoDB *mongo.Database
	clientSingleton sync.Once
)

func EnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    panic("no .env")
  }

  return os.Getenv(key)
}

func MongoDB() *mongo.Database {
	clientSingleton.Do(func () {
		mongoClient = createConnection()
		mongoDB = mongoClient.Database(db)
		
	})

	return mongoDB
}

func createConnection() *mongo.Client {
	// methods creates connection with mongodb
	uri := EnvVariable("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	
	return client
}
