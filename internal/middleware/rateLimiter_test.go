package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

// TestIPRateLimiter_Middleware tests the core rate-limiting functionality.
func TestIPRateLimiter_Middleware(t *testing.T) {
	// A mock handler that will be protected by the middleware.
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a rate limiter that allows 1 request per second with a burst of 1.
	// This makes it easy to test the rate limit.
	limiter := NewIPRateLimiter(rate.Limit(1), 1, time.Second)
	testHandler := limiter.Middleware(mockHandler)

	// --- Test Case 1: First request should be allowed ---
	req1 := httptest.NewRequest("GET", "http://example.com", nil)
	req1.RemoteAddr = "192.0.2.1:12345"
	rr1 := httptest.NewRecorder()

	testHandler.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr1.Code)
	}

	// --- Test Case 2: Second request immediately after should be blocked ---
	req2 := httptest.NewRequest("GET", "http://example.com", nil)
	req2.RemoteAddr = "192.0.2.1:12345" // Same IP
	rr2 := httptest.NewRecorder()

	testHandler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, but got %d", http.StatusTooManyRequests, rr2.Code)
	}

	// --- Test Case 3: A request from a different IP should be allowed ---
	req3 := httptest.NewRequest("GET", "http://example.com", nil)
	req3.RemoteAddr = "198.51.100.2:54321" // Different IP
	rr3 := httptest.NewRecorder()

	testHandler.ServeHTTP(rr3, req3)

	if rr3.Code != http.StatusOK {
		t.Errorf("Expected status code %d for different IP, but got %d", http.StatusOK, rr3.Code)
	}

	// --- Test Case 4: Wait for the rate limit to reset and try again ---
	time.Sleep(1 * time.Second)
	req4 := httptest.NewRequest("GET", "http://example.com", nil)
	req4.RemoteAddr = "192.0.2.1:12345" // Original IP
	rr4 := httptest.NewRecorder()

	testHandler.ServeHTTP(rr4, req4)

	if rr4.Code != http.StatusOK {
		t.Errorf("Expected status code %d after waiting, but got %d", http.StatusOK, rr4.Code)
	}
}

// TestIPRateLimiter_getClientIP tests the logic for extracting the client's IP address.
func TestIPRateLimiter_getClientIP(t *testing.T) {
	limiter := NewIPRateLimiter(1, 1, time.Second)

	testCases := []struct {
		name       string
		header     http.Header
		remoteAddr string
		expectedIP string
	}{
		{
			name:       "X-Forwarded-For (Single IP)",
			header:     http.Header{"X-Forwarded-For": []string{"203.0.113.195"}},
			remoteAddr: "192.0.2.1:12345",
			expectedIP: "203.0.113.195",
		},
		{
			name:       "X-Forwarded-For (Multiple IPs)",
			header:     http.Header{"X-Forwarded-For": []string{"203.0.113.195, 198.51.100.10"}},
			remoteAddr: "192.0.2.1:12345",
			expectedIP: "203.0.113.195",
		},
		{
			name:       "X-Forwarded-For (With spaces)",
			header:     http.Header{"X-Forwarded-For": []string{" 203.0.113.195 , 198.51.100.10"}},
			remoteAddr: "192.0.2.1:12345",
			expectedIP: "203.0.113.195",
		},
		{
			name:       "No X-Forwarded-For Header",
			header:     http.Header{},
			remoteAddr: "192.0.2.1:12345",
			expectedIP: "192.0.2.1:12345",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &http.Request{
				Header:     tc.header,
				RemoteAddr: tc.remoteAddr,
			}
			ip := limiter.getClientIP(req)
			if ip != tc.expectedIP {
				t.Errorf("Expected IP %s, but got %s", tc.expectedIP, ip)
			}
		})
	}
}

// TestIPRateLimiter_getLimiter_Concurrency tests that the getLimiter method is thread-safe.
func TestIPRateLimiter_getLimiter_Concurrency(t *testing.T) {
	limiter := NewIPRateLimiter(10, 10, time.Second)
	ip := "192.168.1.1"

	// Get the initial limiter to compare against.
	initialLimiter := limiter.getLimiter(ip)

	var wg sync.WaitGroup
	numGoroutines := 100

	wg.Add(numGoroutines)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			// Concurrently access the limiter for the same IP.
			l := limiter.getLimiter(ip)
			// Check if the returned limiter is the same instance.
			if l != initialLimiter {
				t.Error("getLimiter returned a different limiter instance for the same IP concurrently")
			}
		}()
	}
	wg.Wait()

	// Check that only one limiter was created for the IP.
	if len(limiter.ips) != 1 {
		t.Errorf("Expected 1 limiter in the map, but found %d", len(limiter.ips))
	}
}
