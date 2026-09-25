package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "modernc.org/sqlite"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	healthchecker "github.com/GuustTaillieu/idiomatic-go/internal/healthchecker"
	"github.com/GuustTaillieu/idiomatic-go/internal/httpapi"
	"github.com/GuustTaillieu/idiomatic-go/internal/httpapi/middleware"
	"github.com/GuustTaillieu/idiomatic-go/internal/processor"
	"github.com/GuustTaillieu/idiomatic-go/internal/store"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db, err := sql.Open("sqlite", "file:orders.db?cache=shared&mode=rwc")
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	orderStore, err := store.NewOrderSqliteStore(db)
	if err != nil {
		slog.Error("Failed to create order store", "error", err)
		os.Exit(1)
	}
	inventoryStore, err := store.NewInventorySQLiteStore(db)
	if err != nil {
		slog.Error("Failed to create inventory store", "error", err)
		os.Exit(1)
	}
	p1 := processor.NewPlaceOrderProcessor(db, inventoryStore, orderStore)
	p2 := processor.NewPaymentProcessor()
	p := processor.NewMultiProcessor(p1, p2)
	queue := domain.NewQueue(p, orderStore, time.Second)
	dispatcher := domain.NewOutboxDispatcher(orderStore, queue)

	healthChecker := healthchecker.NewMultiHealthChecker(store.NewDBHealthChecker(db), queue)

	// Start 3 workers
	queue.Start(ctx, 3)

	// Start the outbox dispatcher
	go func() {
		if err := dispatcher.Start(ctx, 500*time.Millisecond); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("Outbox dispatcher failed", "error", err)
		}
	}()

	// Start HTTP server
	httpHandler := httpapi.NewHandler(orderStore, healthChecker).Routes()
	httpHandler = middleware.RequestIDMiddleware(httpHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 seconds
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server", "error", err)
		}
		queue.Stop()
	}
}
