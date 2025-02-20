package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, _ := zap.NewProduction() // Use zap.NewDevelopment() for a more human-readable format
	defer logger.Sync()              // Ensure logs are flushed to output
	return logger
}

func main() {
	logger := NewLogger()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Failed to load the .env file", zap.Error(err))
	}

	configuredDatastore := os.Getenv("DATASTORE")
	backendURL := ""

	var datastore RateLimiter = getDatastore(configuredDatastore, logger)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := *getClientIP(r)

		// Check the traffic from this IP
		count, err := datastore.Get(clientIP)
		if err != nil {
			logger.Fatal("Unable to get the traffic for IP", zap.String("IP", clientIP))
		}

		if count > 10 {
			http.Error(w, "Too many request", http.StatusTooManyRequests)
		}

		// Store IP in the database
		datastore.Increment(clientIP)

		// Forward the request to the backend
		backendResp, err := http.NewRequest(r.Method, backendURL+r.RequestURI, r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		backendResp.Header = r.Header

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(backendResp)
		if err != nil {
			http.Error(w, "Failed to contact backend", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy the backend response to the client
		for k, v := range resp.Header {
			for _, val := range v {
				w.Header().Add(k, val)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
	logger.Info("Proxy server running on :4000")
	http.ListenAndServe(":4000", nil)

}

func getDatastore(datastore string, logger *zap.Logger) RateLimiter {
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

func getClientIP(r *http.Request) *string {
	clientIP := r.Header.Get("X-Forwarded-For")
	// good when paired with nginx
	if clientIP != "" {
		return &clientIP
	}

	clientIP = r.Header.Get("X-Real-IP")
	if clientIP != "" {
		return &clientIP
	}

	clientIP = r.RemoteAddr
	if clientIP != "" {
		return &clientIP
	}
	return nil
}
