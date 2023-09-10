package db

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The goal here is to always create new version of
// the application as an object if there is any update
func UpsertApplication(ctx context.Context, application Application) (error) {
	coll := MongoDBClient.Database("uptodate").Collection("applications")

	filter := bson.D{{Key: "_id", Value: applicationConvertToId(application.Name, application.Version)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: application.Name},
	{Key: "created_at", Value: application.CreatedAt}, {Key: "version", Value: application.Version},
	{Key: "source", Value: application.Source}, {Key: "vulnerable", Value: application.Vulnerable}}}}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Successfuly updated application: %s", application.Name)
	return nil
}

func GetApplicationByName(ctx context.Context, name string) ([]Application, error) {
	filter := bson.D{{Key: "name", Value: name}}
	sort := bson.D{{Key: "created_at", Value: -1}}
	opts := options.Find().SetSort(sort)
	coll := MongoDBClient.Database("uptodate").Collection("applications")

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []Application
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err

	}
	defer cursor.Close(ctx)
	
	return results, nil
}

// Id is in form of:
// Hash@name%version
// name is lowercase and contains - instead of spaces
// version contains _ instead of dots
// Hashing is for randomness

func applicationConvertToId(name string, version string) string {
	concatenated := name + "%" + version
	resultNoHyphens := strings.ReplaceAll(concatenated, " ", "-")
	resultNoDots := strings.ReplaceAll(resultNoHyphens, ".", "_")
	result := strings.ToLower(resultNoDots)

	hasher := md5.New()
	hasher.Write([]byte(result))
	hash := hex.EncodeToString(hasher.Sum(nil))
	
	id := hash + "@" + result 
	return id
}

func applicationConvertFromId(id string) (string, string) {
	parts := strings.Split(id, "@")

	// Ensure that we have at least two parts
	if len(parts) < 2 {
		fmt.Printf("Invalid id: %s", id)
		return "", ""
	}
	
	firstPart := parts[1]
	subparts := strings.Split(firstPart, "%")

	// Ensure we have two values in subparts
	if len(subparts) < 2 {
		fmt.Printf("Issue parsing id: %s", id)
		return "", ""
	}
	name := subparts[0]
	version := subparts[1]
	name = strings.ReplaceAll(name, "-", " ")
	version = strings.ReplaceAll(version, "_", ".")
	return name, version	
}