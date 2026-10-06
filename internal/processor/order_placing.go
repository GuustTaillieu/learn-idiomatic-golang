package processor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
	"github.com/GuustTaillieu/idiomatic-go/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type OrderPlacing[T any] struct {
	db             *sql.DB
	InventoryStore domain.InventoryStore
	OrderStore     domain.OrderStorer
}

func NewOrderPlacing(db *sql.DB, inventoryStore domain.InventoryStore, orderStore domain.OrderStorer) *OrderPlacing[domain.Order] {
	return &OrderPlacing[domain.Order]{
		db:             db,
		InventoryStore: inventoryStore,
		OrderStore:     orderStore,
	}
}

func (p *OrderPlacing[T]) Process(ctx context.Context, order *domain.Order) (func(context.Context) error, error) {
	ctx = telemetry.ExtractTraceContext(ctx, order.TraceParent)
	ctx, span := telemetry.StartSpan(ctx, "processor.OrderPlacing",
		trace.WithAttributes(
			attribute.String("order.id", order.ID.String()),
		))
	defer span.End()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer tx.Rollback() // Ensure rollback in case of panic or error

	txCtx := lib.WithTx(ctx, tx)

	// Deduct stock
	stock := domain.NewStock(order.ItemID, order.Amount)
	if err := p.InventoryStore.ReserveStock(txCtx, stock); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Update order status
	order.Status = domain.OrderStatusCompleted
	if err := p.OrderStore.Save(txCtx, order); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	if err := tx.Commit(); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to commit order transaction: %w", err)
	}
	// Return a rollback function
	return func(ctx context.Context) error {
		_ = p.InventoryStore.ReleaseStock(ctx, stock)
		order.Status = domain.OrderStatusFailed
		return p.OrderStore.Save(ctx, order)
	}, nil
}
