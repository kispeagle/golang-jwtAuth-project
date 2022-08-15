package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = new()

func new() *mongo.Client {

	dns := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(dns))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func Collection(c *mongo.Client, name string) *mongo.Collection {
	return c.Database("jwt-authentication").Collection(name)
}
