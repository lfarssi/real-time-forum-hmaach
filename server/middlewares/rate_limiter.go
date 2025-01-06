package middlewares

import (
	"net/http"
	"sync"
	"time"
)

type ClientRate struct {
	Requests    int       
	LastRequest time.Time 
}

type RateLimiter struct {
	clients map[string]*ClientRate
	mu      sync.Mutex
	limit   int           
	window  time.Duration // Time window for the rate limit
}


func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*ClientRate),
		limit:   limit,
		window:  window,
	}
}

// isAllowed checks if the client is within the rate limit.
func (rl *RateLimiter) isAllowed(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Check if the client already exists
	if client, exists := rl.clients[clientIP]; exists {
		if now.Sub(client.LastRequest) > rl.window {
			client.Requests = 1
			client.LastRequest = now
			return true
		}

		// Check if the client exceeded the limit
		if client.Requests >= rl.limit {
			return false
		}

		// Increment request count
		client.Requests++
		client.LastRequest = now
		return true
	}

	// Create a new client entry
	rl.clients[clientIP] = &ClientRate{
		Requests:    1,
		LastRequest: now,
	}
	return true
}

// Middleware applies rate limiting to an HTTP handler.
func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr // Extract the client IP
		if !rl.isAllowed(clientIP) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}
