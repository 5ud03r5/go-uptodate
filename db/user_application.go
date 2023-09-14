package db

import "context"

const separator = "#"

func CreateUserApplicationBinding(ctx context.Context, username string, applicationName string) (error) {
	coll := MongoDBClient.Database("uptodate").Collection("user_application")
	userId := UserConvertToId(username)
	appId := ApplicationNameConvertToId(applicationName)
	doc := UserApplication{
		ID: userId + separator + appId,
		ApplicationId: appId,
		UserId: userId,
	}
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

