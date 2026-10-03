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
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusRunning    OrderStatus = "RUNNING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusFailed     OrderStatus = "FAILED"
	OrderStatusDeadLetter OrderStatus = "DEAD_LETTER"
)

type Order struct {
	ID         OrderID
	ItemID     ItemID
	Amount     int
	Status     OrderStatus
	Retries    int
	MaxRetries int
	CreatedAt  time.Time
}

type OrderStorer interface {
	Save(ctx context.Context, order *Order) error
	Get(ctx context.Context, id OrderID) (*Order, error)
	GetPendingOrders(ctx context.Context) ([]*Order, error)
}

type OrderOption func(*Order)

func WithMaxRetries(maxRetries int) OrderOption {
	return func(o *Order) {
		o.MaxRetries = maxRetries
	}
}

func WithOrderStatus(status OrderStatus) OrderOption {
	return func(o *Order) {
		o.Status = status
	}
}

type OrderID uuid.UUID

func NewOrder(itemID ItemID, amount int, opts ...OrderOption) *Order {
	o := &Order{
		ID:         OrderID(uuid.New()),
		ItemID:     itemID,
		Amount:     amount,
		Status:     OrderStatusPending,
		Retries:    0,
		MaxRetries: 3,
		CreatedAt:  time.Now(),
	}
	for _, fn := range opts {
		fn(o)
	}
	return o
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    o.ID.String(),
		Alias: (*Alias)(o),
	})
}

var NilOrderID = OrderID(uuid.Nil())

// MarshalText implements encoding.TextMarshaler
func (id OrderID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler
func (id *OrderID) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalText(data)
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
		if len(v) == 16 {
			*id = OrderID(v)
			return nil
		}
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
