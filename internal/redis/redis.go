package redis

import (
	"os"

	"github.com/go-redis/redis"
)

// New is
func New() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"),
	})

	return redisClient
}
