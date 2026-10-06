package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

type OrderStore struct {
	db lib.DBTX
}

func NewOrderStore(db *sql.DB) (*OrderStore, error) {
	return &OrderStore{db: db}, nil
}

func (s *OrderStore) getDB(ctx context.Context) lib.DBTX {
	if tx, ok := lib.TxFromContext(ctx); ok {
		return tx
	}
	return s.db
}

func (s *OrderStore) Save(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, item_id, amount, status, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			item_id = excluded.item_id,
			amount = excluded.amount,
			status = excluded.status;`
	_, err := s.getDB(ctx).ExecContext(ctx, query, order.ID, order.ItemID, order.Amount, order.Status, order.CreatedAt)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

func (s *OrderStore) Get(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
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
		return nil, domain.ErrInternal
	}
	return &order, nil
}

func (s *OrderStore) GetPendingOrders(ctx context.Context) ([]*domain.Order, error) {
	// Query to get pending orders, limited to 50, ordered by created_at ascending
	query := `
		SELECT id, item_id, amount, status, created_at
		FROM orders
		WHERE status = ?
		ORDER BY created_at ASC
		LIMIT 50;`
	rows, err := s.getDB(ctx).QueryContext(ctx, query, domain.OrderStatusPending)
	if err != nil {
		return nil, domain.ErrInternal
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.ItemID, &order.Amount, &order.Status, &order.CreatedAt); err != nil {
			return nil, domain.ErrInternal
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, domain.ErrInternal
	}
	return orders, nil
}
