package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/time/rate"
	_ "modernc.org/sqlite"

	"github.com/GuustTaillieu/idiomatic-go/internal/config"
	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/event"
	"github.com/GuustTaillieu/idiomatic-go/internal/health"
	httpapi "github.com/GuustTaillieu/idiomatic-go/internal/http"
	"github.com/GuustTaillieu/idiomatic-go/internal/http/middleware"
	"github.com/GuustTaillieu/idiomatic-go/internal/processor"
	"github.com/GuustTaillieu/idiomatic-go/internal/queue"
	"github.com/GuustTaillieu/idiomatic-go/internal/sqlite"
	database "github.com/GuustTaillieu/idiomatic-go/internal/sqlite"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.Load(os.Args[1:]...)
	if err != nil {
		slog.Error("Configuration invalid, failing fast", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	sqlite.Migrate(db)

	orderStore, err := database.NewOrderStore(db)
	if err != nil {
		slog.Error("Failed to create order store", "error", err)
		os.Exit(1)
	}
	inventoryStore, err := database.NewInventoryStore(db)
	if err != nil {
		slog.Error("Failed to create inventory store", "error", err)
		os.Exit(1)
	}
	p1 := processor.NewOrderPlacing(db, inventoryStore, orderStore)
	p2 := processor.NewPaying[*domain.Order]()
	p := processor.NewParallel(p1, p2)
	orderHub := event.NewHub[queue.Task[*domain.Order]]()
	queue := queue.New(p, 100, queue.WithTaskBroadcaster(orderHub))
	dispatcher := domain.NewOutboxDispatcher(orderStore, queue)

	healthChecker := health.NewMultiChecker(health.NewDatabaseChecker(db), queue)

	// Start 3 workers
	queue.Start(ctx, cfg.WorkerCount)

	// Start the outbox dispatcher
	go func() {
		if err := dispatcher.Start(ctx, 500*time.Millisecond); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("Outbox dispatcher failed", "error", err)
		}
	}()

	// Setup rate limiter
	rateLimiter, cleanup := middleware.NewIPRateLimiter(time.Second, rate.Limit(cfg.RateLimit), cfg.RateBurst)
	defer cleanup()

	// Start HTTP server
	httpHandler := httpapi.NewHandler(orderStore, healthChecker, orderHub).Routes()
	httpHandler = middleware.RateLimitMiddleware(rateLimiter)(httpHandler)
	httpHandler = middleware.RequestIDMiddleware(httpHandler)
	httpHandler = middleware.LoggingMiddleware(httpHandler)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: httpHandler,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTime)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Warn("Graceful shutdown timed out, closing forcefully", "error", err)
			_ = srv.Close()
		}
		queue.Stop()
	}
}
