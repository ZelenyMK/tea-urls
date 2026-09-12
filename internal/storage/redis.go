package storage

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		DB:       0,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}
