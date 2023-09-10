package db

import "time"

// Application id needs to be:
// 1. hash(based on name connected with underscores) +
// 2. name (lowercase connected with underscores) +
// 3. version (numbers + lowercase chars)
// So we have them versioned

type Application struct {
	_id string `bson:"_id"`
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

type UserApplication struct {
	_id string `bson:"_id"`
	ApplicationId string `bson:"application_name"`
	UserId string `bson:"user_id"`
}