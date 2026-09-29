package middleware

import (
	"log/slog"
	"net/http"
	"uuid"

	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)

		// Enrich the logger for this entire HTTP request
		logger := slog.With("request_id", requestID)
		ctx := lib.ContextWithLogger(r.Context(), logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
