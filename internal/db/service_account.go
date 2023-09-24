package db

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	auth "github.com/5ud03r5/uptodate/internal/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func CreateServiceAccount(ctx context.Context, applicationName string) (ServiceAccount, error) {
	coll := MongoDBClient.Database("uptodate").Collection("service_accounts")

	password, err1 := auth.GeneratePassword(16)
	if err1 != nil {

		return ServiceAccount{}, err1
	}

	hashedPassword, err2 := auth.HashPassword(password)
	if err2 != nil {
		return ServiceAccount{}, err2
	}

	id := ServiceAccountConvertToId(applicationName)
	doc := ServiceAccount{
		ID: id,
		AccountName: applicationName,
		Password: hashedPassword,
		CreatedAt: time.Now().UTC(),
	}
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return ServiceAccount{}, err
	}
	serviceAccount := doc
	serviceAccount.Password = password
	return serviceAccount, nil
}

func LoginServiceAccount(ctx context.Context, username string, password string) (ServiceAccount, error) {
	id := ServiceAccountConvertToId(username)
	serviceAccount, err := GetServiceAccountByUsername(ctx, id)
	if err != nil {
		return ServiceAccount{}, err
	}
	isPasswordOK := auth.VerifyPassword(password, serviceAccount.Password)
	if !isPasswordOK {
		return ServiceAccount{}, errors.New("password is incorrect")
	}
	return serviceAccount, nil
}

func GetServiceAccountByUsername(ctx context.Context, accountId string) (ServiceAccount, error) {
	filter := bson.D{{Key: "_id", Value: accountId}}
	opts := options.FindOne()
	coll := MongoDBClient.Database("uptodate").Collection("service_accounts")
	var result ServiceAccount
	err := coll.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}


func ServiceAccountConvertToId(username string) string {
	result := strings.ToLower(username)

	hasher := md5.New()
	hasher.Write([]byte(result))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}