package http

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"alive"}`))
}

func (h *Handler) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.healthChecker.Ping(ctx); err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}
