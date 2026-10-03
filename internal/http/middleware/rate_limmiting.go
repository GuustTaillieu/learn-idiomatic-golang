package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func RateLimitMiddleware(limiter *IPRateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if !limiter.GetLimiter(ip).Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type IPRateLimiter struct {
	connections map[string]*connection
	mu          sync.RWMutex

	stopChan    chan struct{}
	cleanupTime time.Duration
	rate        rate.Limit
	burst       int
}
type connection struct {
	rl       *rate.Limiter
	lastSeen time.Time
}

func NewIPRateLimiter(cleanupTime time.Duration, rate rate.Limit, burst int) (rateLimiter *IPRateLimiter, cleanupFunc func()) {
	rateLimiter = &IPRateLimiter{
		connections: make(map[string]*connection),
		mu:          sync.RWMutex{},
		stopChan:    make(chan struct{}),
		cleanupTime: cleanupTime,
		rate:        rate,
		burst:       burst,
	}

	go rateLimiter.startCleaner()

	cleanupFunc = func() {
		close(rateLimiter.stopChan)
	}

	return rateLimiter, cleanupFunc
}

func (l *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	conn, exists := l.connections[ip]
	if !exists {
		conn = &connection{
			rl:       rate.NewLimiter(l.rate, l.burst),
			lastSeen: time.Now(),
		}
		l.connections[ip] = conn
	}
	conn.lastSeen = time.Now()

	return conn.rl
}

func (l *IPRateLimiter) startCleaner() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.mu.Lock()
			for ip, conn := range l.connections {
				if time.Since(conn.lastSeen) > l.cleanupTime {
					delete(l.connections, ip)
				}
			}
			l.mu.Unlock()
		case <-l.stopChan:
			return
		}
	}
}
