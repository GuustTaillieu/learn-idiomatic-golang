package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"uuid"
)

var ErrItemNotFound = errors.New("item not found")

type ItemID uuid.UUID

type Item struct {
	ID        ItemID
	Name      string
	CreatedAt time.Time
}

type ItemOption func(*Item)

func WithID(id ItemID) ItemOption {
	return func(i *Item) {
		i.ID = id
	}
}

func WithCreatedAt(createdAt time.Time) ItemOption {
	return func(i *Item) {
		i.CreatedAt = createdAt
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

var NilItemID = ItemID(uuid.Nil())

// MarshalJSON implements json.Marshaler
func (id ItemID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id))
}

// UnmarshalJSON implements json.Unmarshaler
func (id *ItemID) UnmarshalJSON(data []byte) error {
	uuidPointer := (*uuid.UUID)(id)
	return json.Unmarshal(data, uuidPointer)
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
