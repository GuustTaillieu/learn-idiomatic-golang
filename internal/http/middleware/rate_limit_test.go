package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/http/middleware"
)

// Verify that sending requests up to the burst limit succeeds (200 OK)
func TestRateLimitMiddleware_AllowsRequestsUpToBurstLimit(t *testing.T) {
	limiter, stopFunc := middleware.NewIPRateLimiter(1*time.Second, 1, 5)
	defer stopFunc()

	handler := middleware.RateLimitMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "10.0.0.0"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", rr.Code)
		}
	}
}

// Verify that exceeding the burst limit returns 429 Too Many Requests
func TestRateLimitMiddleware_ExceedsBurstLimit(t *testing.T) {
	limiter, stopFunc := middleware.NewIPRateLimiter(1*time.Second, 1, 5)
	defer stopFunc()

	handler := middleware.RateLimitMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 6; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "10.0.0.0"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if i < 5 && rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK for request %d, got %d", i+1, rr.Code)
		}

		if i == 5 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status 429 Too Many Requests for request %d, got %d", i+1, rr.Code)
		}
	}
}
