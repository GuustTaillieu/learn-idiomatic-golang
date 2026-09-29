package http

import (
	"net/http"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/health"
)

type Handler struct {
	store         domain.OrderStore
	healthChecker health.Checker
}

func NewHandler(store domain.OrderStore, healthChecker health.Checker) *Handler {
	return &Handler{store, healthChecker}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", h.handleCreateorder)
	mux.HandleFunc("GET /orders/{id}", h.handleGetorder)
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("GET /readyz", h.handleReady)

	return mux
}
