package db

import "time"

// Application id needs to be:
// 1. hash(based on name connected with underscores) +
// 2. name (lowercase connected with underscores) +
// 3. version (numbers + lowercase chars)
// So we have them versioned

type RegisteredApplication struct {
	ID string `bson:"_id"`
	Name string `bson:"name"`
	RegisteredAt time.Time `bson:"registered_at"`
	RegisteredBy string `bson:"registered_by"`
}

type Application struct {
	ID string `bson:"_id"`
	Name string	`bson:"name"`
	CreatedAt time.Time	`bson:"created_at"`
	Version string	`bson:"version"`
	Source string	`bson:"source"`
	Vulnerable bool `bson:"vulnerable"`
}

type User struct {
 	ID string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Endpoint string `bson:"endpoint"`
	Email string `'bson':"email"`
	CreatedAt time.Time	`bson:"created_at"`
}

type ServiceAccount struct {
	ID string `bson:"_id,omitempty"`
	AccountName string `bson:"account_name"`
	Password string `bson:"password"`
	CreatedAt time.Time	`bson:"created_at"`
}

type UserApplication struct {
	ID string `bson:"_id,omitempty"`
	ApplicationId string `bson:"application_id"`
	UserId string `bson:"user_id"`
}