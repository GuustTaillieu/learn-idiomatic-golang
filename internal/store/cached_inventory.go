package store

import (
	"context"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"golang.org/x/sync/singleflight"
)

type CachedInventory struct {
	store domain.InventoryStore
	sf    singleflight.Group
}

func NewCachedInventory(store domain.InventoryStore) *CachedInventory {
	return &CachedInventory{
		store: store,
	}
}

func (c *CachedInventory) Get(ctx context.Context, itemID domain.ItemID) (*domain.Stock, error) {
	key := fmt.Sprintf("get_stock:%s", itemID)

	result, err, _ := c.sf.Do(key, func() (any, error) {
		return c.store.Get(ctx, itemID)
	})
	if err != nil {
		return nil, err
	}

	stock, ok := result.(*domain.Stock)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", result)
	}

	return stock, nil
}
