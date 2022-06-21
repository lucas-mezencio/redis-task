package database

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	db   *redis.Client
	once sync.Once
)

func createConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func ConnectWithDB() *redis.Client {
	if db == nil {
		once.Do(func() {
			db = createConnection()
		})
	}
	return db
}
