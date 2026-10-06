package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

type inventoryStore struct {
	db lib.DBTX
	sf lib.Singleflight[*domain.Stock]
}

func (s *inventoryStore) getDB(ctx context.Context) lib.DBTX {
	if tx, ok := lib.TxFromContext(ctx); ok {
		return tx
	}
	return s.db
}

func NewInventoryStore(db *sql.DB) (domain.InventoryStore, error) {
	return &inventoryStore{db: db, sf: lib.Singleflight[*domain.Stock]{}}, nil
}

func (s *inventoryStore) ReserveStock(ctx context.Context, stock *domain.Stock) error {
	query := `
		UPDATE inventory
		SET quantity = quantity - ?, updated_at = ?
		WHERE item_id = ? AND quantity >= ?;`
	res, err := s.getDB(ctx).ExecContext(ctx, query, stock.Quantity, stock.UpdatedAt, stock.ItemID, stock.Quantity)
	if err != nil {
		slog.Error("failed to reserve stock", "error", err)
		return domain.ErrInternal
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return domain.ErrInternal
	}
	if rowsAff == 0 {
		return domain.ErrInsufficientStock
	}
	return nil
}

func (s *inventoryStore) ReleaseStock(ctx context.Context, stock *domain.Stock) error {
	query := `
		UPDATE inventory
		SET quantity = quantity + ?, updated_at = ?
		WHERE item_id = ?;`
	if _, err := s.getDB(ctx).ExecContext(ctx, query, stock.Quantity, stock.UpdatedAt, stock.ItemID); err != nil {
		slog.Error("failed to release stock", "error", err)
		return domain.ErrInternal
	}
	return nil
}

func (s *inventoryStore) AddStock(ctx context.Context, stock *domain.Stock) error {
	query := `
		INSERT INTO inventory (id, item_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			quantity = quantity + excluded.quantity,
			updated_at = ?;`
	if _, err := s.getDB(ctx).ExecContext(ctx, query, stock.ID.String(), stock.ItemID.String(), stock.Quantity, stock.CreatedAt, stock.UpdatedAt, stock.UpdatedAt); err != nil {
		return domain.ErrInternal
	}
	return nil
}

func (s *inventoryStore) Get(ctx context.Context, itemID domain.ItemID) (*domain.Stock, error) {
	key := fmt.Sprintf("get_stock:%s", itemID)

	return s.sf.Do(ctx, key, func() (*domain.Stock, error) {
		query := `
		SELECT item_id, quantity, created_at, updated_at
		FROM inventory
		WHERE item_id = ?;`
		row := s.getDB(ctx).QueryRowContext(ctx, query, itemID)
		var stock domain.Stock
		if err := row.Scan(&stock.ItemID, &stock.Quantity, &stock.CreatedAt, &stock.UpdatedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrStockNotFound
			}
			return nil, domain.ErrInternal
		}
		return &stock, nil
	})
}
