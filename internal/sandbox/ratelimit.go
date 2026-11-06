package sandbox

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter is a sliding-window rate limiter keyed by client IP.
type RateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int           // max allowed calls within the window
	window time.Duration // time window
}

// NewRateLimiter creates a rate limiter: each IP is allowed at most limit times within window seconds.
func NewRateLimiter(limit, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		hits:   map[string][]time.Time{},
		limit:  limit,
		window: time.Duration(windowSeconds) * time.Second,
	}
}

// Allow reports whether the given IP is allowed at this moment.
func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.window)
	ts := r.hits[ip]

	kept := ts[:0]
	for _, t := range ts {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= r.limit {
		r.hits[ip] = kept
		return false
	}
	r.hits[ip] = append(kept, now)
	return true
}

// Wrap wraps an HTTP handler, returning 429 when the limit is exceeded.
func (r *RateLimiter) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ip := clientIP(req)
		if !r.Allow(ip) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"rate limit exceeded, slow down"}`))
			return
		}
		next(w, req)
	}
}

// clientIP extracts the client IP from the request (simple extraction accounting for X-Forwarded-For).
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
