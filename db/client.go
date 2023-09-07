package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client

func DBClient(uri string) (error) {
    // Create a MongoDB client
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
    if err != nil {
        return err
    }

    // Ping the server to verify the connection
    err = client.Ping(context.Background(), nil)
    if err != nil {
        return err
    }

	MongoDBClient = client
    return nil
}
