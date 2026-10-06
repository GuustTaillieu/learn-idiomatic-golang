package domain

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
	"uuid"
)

var (
	ErrInsufficientStock = errors.Join(ErrValidation, errors.New("insufficient stock"))
	ErrStockNotFound     = errors.Join(ErrNotFound, errors.New("stock not found"))
	ErrStockInvalid      = errors.Join(ErrValidation, errors.New("stock is invalid"))
)

type Stock struct {
	ID        StockID
	ItemID    ItemID
	Quantity  int
	UpdatedAt time.Time
	CreatedAt time.Time
}

func NewStock(itemID ItemID, quantity int) *Stock {
	return &Stock{
		ID:        StockID(uuid.New()),
		ItemID:    itemID,
		Quantity:  quantity,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

type InventoryStore interface {
	ReserveStock(ctx context.Context, stock *Stock) error
	ReleaseStock(ctx context.Context, stock *Stock) error
	AddStock(ctx context.Context, stock *Stock) error
	Get(ctx context.Context, itemID ItemID) (*Stock, error)
}

type StockID uuid.UUID

var NilStockID = StockID(uuid.Nil())

// MarshalText implements encoding.TextMarshaler
func (id StockID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler
func (id *StockID) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalText(data)
}

// String implements fmt.Stringer
func (id StockID) String() string {
	return uuid.UUID(id).String()
}

// Value implements driver.Valuer
func (id StockID) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

// Scan implements sql.Scanner
func (id *StockID) Scan(src any) error {
	if src == nil {
		*id = NilStockID
		return nil
	}
	switch v := src.(type) {
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = StockID(u)
		return nil
	case []byte:
		if len(v) == 16 {
			*id = StockID(v)
			return nil
		}
		u, err := uuid.Parse(string(v))
		if err != nil {
			return err
		}
		*id = StockID(u)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into StockID", src)
	}
}
