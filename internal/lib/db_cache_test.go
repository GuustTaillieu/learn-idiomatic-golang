package lib_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

func TestDatabaseCacher_Get(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockInventoryStore := &mockInventoryStore{}
	item := domain.NewItem("Test Item")
	mockStock := domain.NewStock(item.ID, 10)

	var wg sync.WaitGroup
	numGoroutines := 10

	startChan := make(chan struct{})

	// Act
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-startChan:
				res, err := mockInventoryStore.Get(ctx, item.ID)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if res.Quantity != mockStock.Quantity {
					t.Errorf("expected quantity %d, but got %d", mockStock.Quantity, res.Quantity)
				}
			}
		}()
	}
	close(startChan)
	wg.Wait()

	// Assert
	if mockInventoryStore.getCount != 1 {
		t.Errorf("expected Get to be called once, but got %d", mockInventoryStore.getCount)
	}
}

type mockInventoryStore struct {
	getCount int
	mu       sync.Mutex
	sf       lib.Singleflight[*domain.Stock]
}

func (m *mockInventoryStore) Get(ctx context.Context, itemID domain.ItemID) (*domain.Stock, error) {
	key := fmt.Sprintf("stock:%s", itemID.String())
	return m.sf.Do(ctx, key, func() (*domain.Stock, error) {
		m.mu.Lock()
		defer m.mu.Unlock()

		m.getCount++
		time.Sleep(10 * time.Millisecond) // Simulate a delay in fetching from the database
		item := domain.NewItem("Test Item")
		return domain.NewStock(item.ID, 10), nil
	})
}
