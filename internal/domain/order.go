package domain

import (
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
	StatusPending    OrderStatus = "PENDING"
	StatusRunning    OrderStatus = "RUNNING"
	StatusCompleted  OrderStatus = "COMPLETED"
	StatusFailed     OrderStatus = "FAILED"
	StatusDeadLetter OrderStatus = "DEAD_LETTER"
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

func NewOrder(itemID ItemID, amount int, opts ...OrderOption) *Order {
	o := &Order{
		ID:         OrderID(uuid.New()),
		ItemID:     itemID,
		Amount:     amount,
		Status:     StatusPending,
		Retries:    0,
		MaxRetries: 3,
		CreatedAt:  time.Now(),
	}
	for _, fn := range opts {
		fn(o)
	}
	return o
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
