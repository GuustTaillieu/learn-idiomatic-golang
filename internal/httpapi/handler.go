package httpapi

import (
	"context"
	"net/http"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type Queue interface {
	Submit(ctx context.Context, order *domain.Order) error
}

type OrderStore interface {
	Get(ctx context.Context, id domain.OrderID) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}

type Handler struct {
	queue Queue
	store OrderStore
}

func NewHandler(queue Queue, store OrderStore) *Handler {
	return &Handler{queue, store}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", h.handleCreateorder)
	mux.HandleFunc("GET /orders/{id}", h.handleGetorder)

	return mux
}
