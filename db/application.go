package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The goal here is to always create new version of
// the application as an object if there is any update
func UpsertApplication(application Application) (error) {
	coll := MongoDBClient.Database("uptodate").Collection("applications")

	filter := bson.D{{Key: "_id", Value: ConvertToId(application.Name, application.Version)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: application.Name},
	{Key: "created_at", Value: application.CreatedAt}, {Key: "version", Value: application.Version},
	{Key: "source", Value: application.Source}, {Key: "vulnerable", Value: application.Vulnerable}}}}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		fmt.Printf("Error while updating application: %s", err)
		return err
	}
	fmt.Printf("Successfuly updated application: %s", application.Name)
	return nil
}

func GetApplicationByName(name string) ([]Application, error) {
	filter := bson.D{{Key: "name", Value: name}}
	sort := bson.D{{Key: "version", Value: 1}}
	opts := options.Find().SetSort(sort)
	coll := MongoDBClient.Database("uptodate").Collection("applications")
	cursor, err := coll.Find(context.Background(), filter, opts)
	if err != nil {
		fmt.Printf("Could not find applications: %s", err)
		return nil, err
	}
	var results []Application
	if err = cursor.All(context.Background(), &results); err != nil {
		fmt.Printf("Error getting applications: %s", err)
		return nil, err

	}
	return results, nil
}