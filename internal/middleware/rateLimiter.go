package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*limiterEntry

	mu       sync.Mutex
	r        rate.Limit
	b        int
	lastSeen time.Duration
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewIPRateLimiter is a constructor for IPRateLimiter.
func NewIPRateLimiter(r rate.Limit, b int, lastSeen time.Duration) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:      make(map[string]*limiterEntry),
		r:        r,
		b:        b,
		lastSeen: lastSeen,
	}
	go i.cleanupVisitors()
	return i
}

// getLimiter retrieves or creates a limiter for a given IP address.
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	// Lock the mutex to protect the map.
	i.mu.Lock()
	defer i.mu.Unlock()

	entry, exists := i.ips[ip]
	// If the IP is not in the map, create a new limiter and add it.
	if !exists {
		limiter := rate.NewLimiter(i.r, i.b)
		i.ips[ip] = &limiterEntry{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

func (i *IPRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(10 * time.Minute)
		i.mu.Lock()
		for ip, entry := range i.ips {
			if time.Since(entry.lastSeen) > i.lastSeen {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}

func (i *IPRateLimiter) getClientIP(r *http.Request) string {
	// Get the X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")

	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// If X-Forwarded-For is not present, fall back to RemoteAddr
	return r.RemoteAddr
}
func (i *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := i.getClientIP(r)
		limiter := i.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)

	})
}
