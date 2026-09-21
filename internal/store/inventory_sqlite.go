package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

type InventorySQLite struct {
	db lib.DBTX
}

func (s *InventorySQLite) getDB(ctx context.Context) lib.DBTX {
	if tx, ok := lib.TxFromContext(ctx); ok {
		return tx
	}
	return s.db
}

func NewInventorySQLiteStore(db *sql.DB) (*InventorySQLite, error) {
	query := `
		CREATE TABLE IF NOT EXISTS inventory (
			id TEXT PRIMARY KEY,
			item_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
		);`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	return &InventorySQLite{db: db}, nil
}

func (s *InventorySQLite) ReserveStock(ctx context.Context, stock *domain.Stock) error {
	query := `
		UPDATE inventory
		SET quantity = quantity - ?, updated_at = ?
		WHERE item_id = ? AND quantity >= ?;`
	res, err := s.getDB(ctx).ExecContext(ctx, query, stock.Quantity, stock.UpdatedAt, stock.ItemID, stock.Quantity)
	if err != nil {
		return err
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

func (s *InventorySQLite) ReleaseStock(ctx context.Context, stock *domain.Stock) error {
	query := `
		UPDATE inventory
		SET quantity = quantity + ?, updated_at = ?
		WHERE item_id = ?;`
	if _, err := s.getDB(ctx).ExecContext(ctx, query, stock.Quantity, stock.UpdatedAt, stock.ItemID); err != nil {
		return err
	}
	return nil
}

func (s *InventorySQLite) AddStock(ctx context.Context, stock *domain.Stock) error {
	query := `
		INSERT INTO inventory (id, item_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			quantity = quantity + excluded.quantity,
			updated_at = ?;`
	if _, err := s.getDB(ctx).ExecContext(ctx, query, stock.ID.String(), stock.ItemID.String(), stock.Quantity, stock.CreatedAt, stock.UpdatedAt, stock.UpdatedAt); err != nil {
		return fmt.Errorf("failed to add stock: %w", err)
	}
	return nil
}

func (s *InventorySQLite) Get(ctx context.Context, itemID domain.ItemID) (*domain.Stock, error) {
	query := `
		SELECT item_id, quantity, created_at, updated_at
		FROM inventory
		WHERE item_id = ?;`
	row := s.getDB(ctx).QueryRowContext(ctx, query, itemID)
	var stock domain.Stock
	if err := row.Scan(&stock.ItemID, &stock.Quantity, &stock.CreatedAt, &stock.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("stock not found for item_id %s: %w", itemID, err)
		}
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}
	return &stock, nil
}
