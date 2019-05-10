package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var server *mongo.Client // Create Global Variable for Share library collection in this package

//Connect function for Connect database and get collections
func Connect(dburl string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // Create Context for connect to server. This connection has 10 minute timeout for connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl)) // Connect Database
	if err != nil {
		log.Fatal("error : " + err.Error())
	}
	server = client
}

func GetDB(name string) *mongo.Collection {
	collection := server.Database("library").Collection(name) //Getting library collection from database
	return collection
}
