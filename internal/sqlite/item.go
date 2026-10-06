package sqlite

import (
	"context"
	"database/sql"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type itemStore struct {
	db *sql.DB
}

func NewItemStore(db *sql.DB) (domain.ItemStorer, error) {
	return &itemStore{db: db}, nil
}

func (s *itemStore) Save(ctx context.Context, item *domain.Item) error {
	query := `
		INSERT INTO items (id, name)
		VALUES (?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name;`
	_, err := s.db.ExecContext(ctx, query, item.ID, item.Name)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

func (s *itemStore) Get(ctx context.Context, id domain.ItemID) (*domain.Item, error) {
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
		return nil, domain.ErrInternal
	}
	return &item, nil
}
