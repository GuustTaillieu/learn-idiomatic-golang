package httpapi

import (
	"context"
	"net/http"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type Queue interface {
	Submit(ctx context.Context, order *domain.Order) error
}

type OrderStore interface {
	Get(ctx context.Context, id domain.OrderID) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}

type Handler struct {
	store         OrderStore
	healthChecker HealthChecker
}

func NewHandler(store OrderStore, healthChecker HealthChecker) *Handler {
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
