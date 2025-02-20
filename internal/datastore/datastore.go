package datastore

import (
	"go.uber.org/zap"
)

type RateLimiter interface {
	Increment(key string) (int, error)
	Get(key string) (int, error)
	Set(key string, value int) error
	Reset(key string) error
}

func NewDatastore(datastore string, logger *zap.Logger) RateLimiter {
	switch datastore {
	case "sqlite":
		rl, err := NewSQLiteRateLimiter("rate_limits.db", logger)
		if err != nil {
			logger.Fatal("Failed to connect to SQLite", zap.Error(err))
		}
		return rl
	case "redis":
		return NewRedisRateLimiter("localhost:6379", logger)
	default:
		logger.Fatal("Unsupported datastore")
	}

	return nil
}
