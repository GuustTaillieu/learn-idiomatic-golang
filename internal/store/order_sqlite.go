package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

type OrderSqliteStore struct {
	db lib.DBTX
}

func NewOrderSqliteStore(db *sql.DB) (*OrderSqliteStore, error) {
	query := `
		CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			item_id TEXT NOT NULL,
			amount INTEGER NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY(item_id) REFERENCES items(id)
		);`
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create orders table: %w", err)
	}
	return &OrderSqliteStore{db: db}, nil
}

func (s *OrderSqliteStore) getDB(ctx context.Context) lib.DBTX {
	if tx, ok := lib.TxFromContext(ctx); ok {
		return tx
	}
	return s.db
}

func (s *OrderSqliteStore) Save(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, item_id, amount, status, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			item_id = excluded.item_id,
			amount = excluded.amount,
			status = excluded.status;`
	_, err := s.getDB(ctx).ExecContext(ctx, query, order.ID, order.ItemID, order.Amount, order.Status, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	return nil
}

func (s *OrderSqliteStore) Get(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	query := `
		SELECT id, item_id, amount, status, created_at
		FROM orders
		WHERE id = ?;`
	row := s.getDB(ctx).QueryRowContext(ctx, query, id)
	var order domain.Order
	if err := row.Scan(&order.ID, &order.ItemID, &order.Amount, &order.Status, &order.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return &order, nil
}

func (s *OrderSqliteStore) GetPendingOrders(ctx context.Context) ([]*domain.Order, error) {
	query := `
		SELECT id, item_id, amount, status, created_at
		FROM orders
		WHERE status = ?
		ORDER BY created_at ASC;`
	rows, err := s.getDB(ctx).QueryContext(ctx, query, domain.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.ItemID, &order.Amount, &order.Status, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over orders: %w", err)
	}
	return orders, nil
}
