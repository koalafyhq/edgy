package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

var edgyRegion = os.Getenv("EDGY_REGION")

// Error is
func Error(reqID string, err error, msg string) {
	log.Error().Err(err).Str("req_id", reqID).Str("edge", edgyRegion).Msg(msg)
}

// Fatal is
func Fatal(reqID string, err error, msg string) {
	log.Fatal().Err(err).Str("req_id", reqID).Str("edge", edgyRegion).Msg(msg)
}

// Debug is
func Debug(msg string) {
	log.Debug().Str("edge", edgyRegion).Msg(msg)
}

// Cache is
func Cache(status string, endpoint string, path string) {
	log.Info().Str("edge", edgyRegion).Str("cache", status).Str("endpoint", endpoint).Str("path", path).Msg("")
}

// CacheWithSize is
func CacheWithSize(status string, endpoint string, path string, size int64) {
	log.Info().Str("edge", edgyRegion).Str("cache", status).Str("endpoint", endpoint).Str("path", path).Int64("size", size).Msg("")
}

// CacheWithLatency is
func CacheWithLatency(status string, endpoint string, path string, latency time.Duration) {
	log.Info().Str("edge", edgyRegion).Str("cache", status).Dur("latency", latency).Str("path", path).Msg("")
}
