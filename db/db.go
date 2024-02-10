package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	dbName = "contacts"
)

var dbUrl string = os.Getenv("MONGO_DB_URL")

// initialize mongodb connection
func DbConnection() {
	clientOptions := options.Client().ApplyURI(dbUrl)

	//connect to db
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		panic(err)
	}
 
	//check connection
	err = client.Ping(ctx, nil) 
	if err != nil {
		panic(err)
	} 
	log.Println("Connected to MongoDB")
}

// mongo client instance
func GetClient() *mongo.Client {
	return Client
}

// get collection instance
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database(dbName).Collection(collectionName)
}

// diconnect the mongodb connection
func DbDisconnect() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := Client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
		log.Println("Disconnected from MongoDB!")
	}

}
