package datastore

import (
	"database/sql"

	"go.uber.org/zap"
)

type SQLiteRateLimiter struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSQLiteRateLimiter(dsn string, logger *zap.Logger) (*SQLiteRateLimiter, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		logger.Error("Failed to open the sqlite datastore")
		return nil, err
	}
	return &SQLiteRateLimiter{db: db, logger: logger}, nil
}

func (r *SQLiteRateLimiter) Increment(key string) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT count FROM rate_limits WHERE key = ?", key).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Error("Failed to open the sqlite datastore")
		return 0, err
	}
	if err == sql.ErrNoRows {
		count = 0
	}

	count++
	_, err = r.db.Exec("INSERT INTO rate_limits (key, count) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET count = ?", key, count, count)
	return count, err
}

func (r *SQLiteRateLimiter) Get(key string) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT count FROM rate_limits WHERE key = ?", key).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *SQLiteRateLimiter) Set(key string, value int) error {
	_, err := r.db.Exec("INSERT INTO rate_limits (key, count) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET count = ?", key, value, value)
	return err
}

func (r *SQLiteRateLimiter) Reset(key string) error {
	_, err := r.db.Exec("DELETE FROM rate_limits WHERE key = ?", key)
	return err
}
