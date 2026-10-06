package domain

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
	"uuid"
)

var ErrItemNotFound = errors.Join(ErrNotFound, errors.New("item not found"))

type Item struct {
	ID        ItemID
	Name      string
	CreatedAt time.Time
}

type ItemStorer interface {
	Save(ctx context.Context, item *Item) error
	Get(ctx context.Context, id ItemID) (*Item, error)
}

type ItemOption func(*Item)

func WithItemID(id ItemID) ItemOption {
	return func(i *Item) {
		i.ID = id
	}
}

func NewItem(name string, opts ...ItemOption) *Item {
	i := &Item{
		ID:        ItemID(uuid.New()),
		Name:      name,
		CreatedAt: time.Now(),
	}
	for _, fn := range opts {
		fn(i)
	}
	return i
}

type ItemID uuid.UUID

var NilItemID = ItemID(uuid.Nil())

// MarshalText implements encoding.TextMarshaler
func (id ItemID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler
func (id *ItemID) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalText(data)
}

// String implements fmt.Stringer
func (id ItemID) String() string {
	return uuid.UUID(id).String()
}

// Value implements driver.Valuer
func (id ItemID) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

// Scan implements sql.Scanner
func (id *ItemID) Scan(src any) error {
	if src == nil {
		*id = NilItemID
		return nil
	}
	switch v := src.(type) {
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = ItemID(u)
		return nil
	case []byte:
		if len(v) == 16 {
			*id = ItemID(v)
			return nil
		}
		u, err := uuid.Parse(string(v))
		if err != nil {
			return err
		}
		*id = ItemID(u)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into ItemID", src)
	}
}
