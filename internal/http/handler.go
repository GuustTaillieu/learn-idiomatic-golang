package http

import (
	"net/http"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/health"
	"github.com/GuustTaillieu/idiomatic-go/internal/queue"
	"github.com/GuustTaillieu/idiomatic-go/web"
)

type Handler struct {
	orderStore    domain.OrderStorer
	healthChecker health.Checker
	orderHub      domain.Hub[queue.Task[*domain.Order]]
}

func NewHandler(store domain.OrderStorer, healthChecker health.Checker, orderHub domain.Hub[queue.Task[*domain.Order]]) *Handler {
	return &Handler{
		orderStore:    store,
		healthChecker: healthChecker,
		orderHub:      orderHub,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", h.handleCreateorder)
	mux.HandleFunc("GET /orders/{id}", h.handleGetorder)

	mux.HandleFunc("GET /events", h.handleGetEvents)

	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("GET /readyz", h.handleReady)

	mux.Handle("GET /", web.FileServer())

	return mux
}
