package middleware

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter

	mu sync.Mutex
	r  rate.Limit
	b  int
}

// NewIPRateLimiter is a constructor for IPRateLimiter.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

// getLimiter retrieves or creates a limiter for a given IP address.
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	// Lock the mutex to protect the map.
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	// If the IP is not in the map, create a new limiter and add it.
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
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
