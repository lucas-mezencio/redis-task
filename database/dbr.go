// Package database is responsible for the database connections management
package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	db   *redis.Client
	once sync.Once
)

// CTX global variable for the context of database functions
var CTX = context.Background()

// createConnection create the database connection
func createConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// ConnectWithDB creates a singleton variable for the database connection
// Is based on createConnection hidden function
func ConnectWithDB() *redis.Client {
	if db == nil {
		once.Do(func() {
			db = createConnection()
		})
	}
	return db
}
