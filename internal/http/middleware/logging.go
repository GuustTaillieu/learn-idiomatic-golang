package middleware

import (
	"net/http"
	"time"

	httpapi "github.com/GuustTaillieu/idiomatic-go/internal/http"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := httpapi.NewResponseRecorder(w)

		next.ServeHTTP(rr, r)

		lib.Logger(r.Context()).Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rr.StatusCode,
			"bytes", rr.Bytes,
			"duration", time.Since(start),
		)
	})
}
