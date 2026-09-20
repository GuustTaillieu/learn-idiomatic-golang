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
	"github.com/GuustTaillieu/idiomatic-go/internal/httpapi"
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
	q := domain.NewQueue(p, orderStore)

	// Start 3 workers
	q.Start(ctx, 3)

	// Before starting the HTTP server, we process any pending orders
	pendingOrders, err := orderStore.GetPendingOrders(ctx)
	if err != nil {
		slog.Error("Failed to get pending orders", "error", err)
		os.Exit(1)
	}
	for _, order := range pendingOrders {
		if err := q.Submit(ctx, order); err != nil {
			slog.Error("Failed to submit pending order", "order_id", order.ID, "error", err)
		}
	}

	// Start HTTP server
	handler := httpapi.NewHandler(q, orderStore)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
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
		q.Stop()
	}
}
