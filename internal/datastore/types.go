package datastore

// RateLimiter defines the common methods for both SQLite and Redis backends.
type RateLimiter interface {
	Increment(key string) (int, error)
	Get(key string) (int, error)
	Set(key string, value int) error
	Reset(key string) error
}
