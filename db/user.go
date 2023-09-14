package db

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/5ud03r5/uptodate/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterUser(ctx context.Context, username string, email string, endpoint string) (User, error) {
	coll := MongoDBClient.Database("uptodate").Collection("users")

	password, err1 := auth.GeneratePassword(16)
	if err1 != nil {

		return User{}, err1
	}

	hashedPassword, err2 := auth.HashPassword(password)
	if err2 != nil {
		return User{}, err2
	}

	id := UserConvertToId(username)
	doc := User{
		ID: id,
		Username: username,
		Password: hashedPassword,
		Email: email,
		Endpoint: endpoint,
		CreatedAt: time.Now().UTC(),
	}
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return User{}, err
	}
	user := doc
	user.Password = password
	return user, nil
}

func LoginUser(ctx context.Context, username string, password string) (User, error) {
	id := UserConvertToId(username)
	user, err := GetUserByUsername(ctx, id)
	if err != nil {
		return User{}, err
	}
	isPasswordOK := auth.VerifyPassword(password, user.Password)
	if !isPasswordOK {
		return User{}, errors.New("password is incorrect")
	}
	return user, nil
}

func GetUserByUsername(ctx context.Context, username string) (User, error) {
	filter := bson.D{{Key: "_id", Value: username}}
	opts := options.FindOne()
	coll := MongoDBClient.Database("uptodate").Collection("users")
	var result User
	err := coll.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func UserConvertToId(username string) string {
	result := strings.ToLower(username)

	hasher := md5.New()
	hasher.Write([]byte(result))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

func UserConvertFromId(id string) (string) {
	// TODO	
	return ""
}