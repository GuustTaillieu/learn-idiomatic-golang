package http

import (
	"net/http"
	"net/http/pprof"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/health"
	"github.com/GuustTaillieu/idiomatic-go/internal/queue"
	"github.com/GuustTaillieu/idiomatic-go/proto/order/v1/orderv1connect"
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

	osp, osh := orderv1connect.NewOrderServiceHandler(NewOrderRPCServer(h.orderStore))

	mux.HandleFunc("POST /orders", h.handleCreateorder)
	mux.HandleFunc("GET /orders/{id}", h.handleGetorder)
	mux.Handle(osp, osh)

	mux.HandleFunc("GET /events", h.handleGetEvents)

	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("GET /readyz", h.handleReady)
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)

	mux.Handle("/", web.FileServer())

	return mux
}
