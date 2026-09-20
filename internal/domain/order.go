package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"uuid"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusRunning   OrderStatus = "RUNNING"
	StatusCompleted OrderStatus = "COMPLETED"
	StatusFailed    OrderStatus = "FAILED"
)

type Order struct {
	ID        OrderID
	ItemID    ItemID
	Amount    int
	Status    OrderStatus
	CreatedAt time.Time
}

func NewOrder(itemID ItemID, amount int) *Order {
	return &Order{
		ID:        OrderID(uuid.New()),
		ItemID:    itemID,
		Amount:    amount,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}

type OrderID uuid.UUID

var NilOrderID = OrderID(uuid.Nil())

// MarshalJSON implements json.Marshaler
func (id OrderID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id))
}

// UnmarshalJSON implements json.Unmarshaler
func (id *OrderID) UnmarshalJSON(data []byte) error {
	uuidPointer := (*uuid.UUID)(id)
	return json.Unmarshal(data, uuidPointer)
}

// String implements fmt.Stringer
func (id OrderID) String() string {
	return uuid.UUID(id).String()
}

// Value implements driver.Valuer
func (id OrderID) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

// Scan implements sql.Scanner
func (id *OrderID) Scan(src any) error {
	if src == nil {
		*id = NilOrderID
		return nil
	}
	switch v := src.(type) {
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = OrderID(u)
		return nil
	case []byte:
		u, err := uuid.Parse(string(v))
		if err != nil {
			return err
		}
		*id = OrderID(u)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into OrderID", src)
	}
}

type Processor interface {
	Process(ctx context.Context, order *Order) (func() error, error)
}
