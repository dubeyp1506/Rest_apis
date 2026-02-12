package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	resetTime time.Duration
	limit     int
}

func NewRateLimiter(limit int, resetTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}

	go rl.ResetVisitorsCount()

	return rl
}

func (rl *RateLimiter) ResetVisitorsCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		visitorIp := r.RemoteAddr
		if rl.visitors[visitorIp] > rl.limit {
			http.Error(w, "Too Many Request", http.StatusTooManyRequests)
			return
		}
		rl.visitors[visitorIp]++
		fmt.Printf("Visitor IP: %s, Request Count: %d\n", visitorIp, rl.visitors[visitorIp])

		next.ServeHTTP(w, r)
	})
}
