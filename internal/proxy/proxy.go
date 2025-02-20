package proxy

import (
	"github.com/gvarma28/MiniRL/internal/datastore"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Proxy struct {
	backendURL string
	store      datastore.RateLimiter
	logger     *zap.Logger
}

func NewProxy(backendURL string, rateLimiter datastore.RateLimiter, logger *zap.Logger) http.Handler {
	return &Proxy{
		backendURL: backendURL,
		store:      rateLimiter,
		logger:     logger,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := *getClientIP(r)

	// Rate limit check
	count, err := p.store.Get(clientIP)
	if err != nil {
		p.logger.Error("Failed to get rate limit", zap.String("IP", clientIP), zap.Error(err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if count >= 10 {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	p.store.Increment(clientIP)

	// Forward the request
	p.logger.Info("Forwarding request", zap.String("IP", clientIP))
	resp, err := http.NewRequest(r.Method, p.backendURL+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	resp.Header = r.Header
	client := &http.Client{}
	backendResp, err := client.Do(resp)
	if err != nil {
		http.Error(w, "Backend error", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	// Copy response
	for k, v := range backendResp.Header {
		for _, val := range v {
			w.Header().Add(k, val)
		}
	}
	w.WriteHeader(backendResp.StatusCode)
	io.Copy(w, backendResp.Body)
}

func getClientIP(r *http.Request) *string {
	clientIP := r.Header.Get("X-Forwarded-For")

	if clientIP != "" {
		return &clientIP
	}

	// good when paired with nginx
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
