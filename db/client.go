package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client

type CollectionIndex struct {
	CollectionName string
	IndexType string
	IndexField string
}


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

// Create indexes on specified fields in collections
func CreateIndexes(collections ...CollectionIndex) {
	for _, collection := range collections {
		coll := MongoDBClient.Database("uptodate").Collection(collection.CollectionName)
		indexModel := mongo.IndexModel{
			Keys:    bson.M{collection.IndexField: collection.IndexType},
		}
		_, err := coll.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			fmt.Printf("Error creating text index: %v\n", err)
			return
		}
	}
}
