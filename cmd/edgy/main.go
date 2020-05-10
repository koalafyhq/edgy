package main

import (
	"github.com/koalafy/edgy/internal/cachemanager"
	"github.com/koalafy/edgy/internal/logger"
	"github.com/koalafy/edgy/internal/redis"
	"github.com/koalafy/edgy/internal/server"
)

func main() {
	redisClient := redis.New()

	if _, err := redisClient.Ping().Result(); err != nil {
		logger.Fatal("", err, "failed to connect to the redis instance")
	}

	defer redisClient.Close()

	cacheManager := cachemanager.UseRedis(redisClient)
	app := server.New(cacheManager)

	if err := server.Run(app); err != nil {
		logger.Error("", err, "failed to start the server")
	}
}
