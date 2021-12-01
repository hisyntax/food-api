package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("Error finding .env file")
	}
	MongoDb := os.Getenv("MONGODB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Connected to MongoDB!")

	return client
}

//client Database instance
var Client *mongo.Client = DBInstance()

//openCollection is a function that makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	clusterName := os.Getenv("CLUSTER_NAME")

	var collection *mongo.Collection = client.Database(clusterName).Collection(os.Getenv("DB_COLLECTION_NAME"))
	return collection
}
