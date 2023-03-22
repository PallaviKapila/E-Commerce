package database

import (
	"context"
	"fmt"
	"log" //to logout errors
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv" //to work with enviornment variables
	"go.mongodb.org/mongo-driver/mongo"
)

// it returns mongo client
func DBinstance() *mongo.Client {
	//LOAD FUNCTION TO LOAD UP OUUR ENVIORNMENT
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
	//create variable called mongodb
	//os package to get the enviornment for MongoDb
	MongoDb := os.Getenv("MONGODB_URL")
	//here we'll pass the mongodb variable that we created and we'll capture this in variable client & cspture the error as well
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer function to call cancel at the end of this function
	defer cancel()
	//we'll use context to connect to our client
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	//this client is basically mongo client dereference to it which we are returning from this function
	return client
}

// we call this function which runs mongo client which we have captured in this variable client which is of the type mongo client
var Client *mongo.Client = DBinstance()

// to access particular collection in database
// we can pass a particular collection here and then we can use it
// this function returns a particular collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	//we define a variable collection which is of type mongo.Collection it's a client.Database and we can define the database here and we can define the name of the collection
	var collection *mongo.Collection = client.Database("authentication").Collection(collectionName)
	//we'll return the collection here
	return collection
	//so this database can be the database that we get a fully managed database in a mongodb instance in the cloud o the database that we'll get would be mostly like cluster-0 or cluster-1
}
