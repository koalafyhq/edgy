package cachemanager

import "github.com/go-redis/redis"

// CacheManager is
type CacheManager struct {
	redis *redis.Client
}

// UseRedis is
func UseRedis(redis *redis.Client) *CacheManager {
	return &CacheManager{
		redis: redis,
	}
}
