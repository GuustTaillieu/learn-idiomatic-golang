package processor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

type OrderPlacing[T any] struct {
	db             *sql.DB
	InventoryStore domain.InventoryStore
	OrderStore     domain.OrderStore
}

func NewOrderPlacing(db *sql.DB, inventoryStore domain.InventoryStore, orderStore domain.OrderStore) *OrderPlacing[domain.Order] {
	return &OrderPlacing[domain.Order]{
		db:             db,
		InventoryStore: inventoryStore,
		OrderStore:     orderStore,
	}
}

func (p *OrderPlacing[T]) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Ensure rollback in case of panic or error

	txCtx := lib.WithTx(ctx, tx)

	// Deduct stock
	stock := domain.NewStock(order.ItemID, order.Amount)
	if err := p.InventoryStore.ReserveStock(txCtx, stock); err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Insert the order
	if err := p.OrderStore.Save(txCtx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit order transaction: %w", err)
	}
	// Return a rollback function
	return func() error {
		if err := p.InventoryStore.ReleaseStock(ctx, stock); err != nil {
			return err
		}

		return nil
	}, nil
}
