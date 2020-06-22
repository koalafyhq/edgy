package main

import (
	"github.com/rs/zerolog/log"

	"github.com/koalafy/edgy/internal/cachemanager"
	"github.com/koalafy/edgy/internal/redis"
	"github.com/koalafy/edgy/internal/server"
)

func main() {
	redisClient := redis.New()

	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatal().Msg("failed to connect to the redis instance")
	}

	defer redisClient.Close()

	cacheManager := cachemanager.UseRedis(redisClient)
	app := server.New(cacheManager)

	if err := server.Run(app); err != nil {
		log.Error().Err(err).Msg("")
	}
}
