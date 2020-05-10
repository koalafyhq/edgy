package cachemanager

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

const cacheID = "cache:"
const endpointID = "endpoint:"
const oneYear = 8760 * time.Hour

// GetEndpoint is
func (c *CacheManager) GetEndpoint(host string) ([]byte, error) {
	cacheKey := endpointID + host
	value, err := c.redis.Get(cacheKey).Bytes()

	if err == redis.Nil {
		return nil, errors.New(err.Error())
	}

	return value, err
}

// Get is
func (c *CacheManager) Get(key string) ([]byte, error) {
	cacheKey := cacheID + key
	value, err := c.redis.Get(cacheKey).Bytes()

	if err == redis.Nil {
		return nil, errors.New(err.Error())
	}

	return value, err
}

// Set is
func (c *CacheManager) Set(key string, value []byte) {
	cacheKey := cacheID + key

	// set default cache to 1yr
	c.redis.Set(cacheKey, value, oneYear)
}
