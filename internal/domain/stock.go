package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"uuid"
)

var ErrStockNotFound = errors.New("stock not found")

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

type StockID uuid.UUID

var NilStockID = StockID(uuid.Nil())

// MarshalJSON implements json.Marshaler
func (id StockID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id))
}

// UnmarshalJSON implements json.Unmarshaler
func (id *StockID) UnmarshalJSON(data []byte) error {
	uuidPointer := (*uuid.UUID)(id)
	return json.Unmarshal(data, uuidPointer)
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
