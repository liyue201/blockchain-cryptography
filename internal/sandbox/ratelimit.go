package sandbox

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter 基于客户端 IP 的滑动窗口限流。
type RateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int           // 窗口内允许次数
	window time.Duration // 时间窗口
}

// NewRateLimiter 创建限流器：每个 IP 在 window 秒内最多 limit 次。
func NewRateLimiter(limit, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		hits:   map[string][]time.Time{},
		limit:  limit,
		window: time.Duration(windowSeconds) * time.Second,
	}
}

// Allow 判断该 IP 此刻是否放行。
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

// Wrap 包装 HTTP handler，超限时返回 429。
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

// clientIP 从请求提取客户端 IP（考虑 X-Forwarded-For 前的简单取值）。
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
