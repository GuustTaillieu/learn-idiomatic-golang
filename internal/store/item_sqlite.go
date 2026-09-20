package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type ItemSQLiteStore struct {
	db *sql.DB
}

func NewItemSQLiteStore(db *sql.DB) (*ItemSQLiteStore, error) {
	query := `
		CREATE TABLE IF NOT EXISTS items (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	return &ItemSQLiteStore{db: db}, nil
}

func (s *ItemSQLiteStore) Save(ctx context.Context, item *domain.Item) error {
	query := `
		INSERT INTO items (id, name)
		VALUES (?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name;`
	_, err := s.db.ExecContext(ctx, query, item.ID, item.Name)
	if err != nil {
		return fmt.Errorf("failed to save item: %w", err)
	}
	return nil
}

func (s *ItemSQLiteStore) Get(ctx context.Context, id domain.ItemID) (*domain.Item, error) {
	query := `
		SELECT id, name
		FROM items
		WHERE id = ?;`
	row := s.db.QueryRowContext(ctx, query, id)
	var item domain.Item
	if err := row.Scan(&item.ID, &item.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	return &item, nil
}
