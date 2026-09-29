package store_test

import (
	"context"
	"sync"
	"testing"
	"time"
	"uuid"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/store"
)

// Write a test with a mock store that counts how many times Get is executed.
func TestCachedInventory_Get(t *testing.T) {
	// Arrange
	itemID := domain.ItemID(uuid.New())
	mockStock := domain.NewStock(itemID, 10)
	mockStore := NewMockInventoryStore(mockStock)
	cachedInventory := store.NewCachedInventory(mockStore)

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
				res, err := cachedInventory.Get(context.Background(), itemID)
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
	if mockStore.getCount != 1 {
		t.Errorf("expected Get to be called once, but got %d", mockStore.getCount)
	}
}

type mockInventoryStore struct {
	getCount  int
	mu        sync.Mutex
	mockStock *domain.Stock
}

func NewMockInventoryStore(stock *domain.Stock) *mockInventoryStore {
	return &mockInventoryStore{
		mockStock: stock,
	}
}

func (m *mockInventoryStore) ReserveStock(ctx context.Context, stock *domain.Stock) error {
	return nil
}

func (m *mockInventoryStore) ReleaseStock(ctx context.Context, stock *domain.Stock) error {
	return nil
}

func (s *mockInventoryStore) AddStock(ctx context.Context, stock *domain.Stock) error {
	return nil
}

func (m *mockInventoryStore) Get(ctx context.Context, itemID domain.ItemID) (*domain.Stock, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.getCount++
	time.Sleep(10 * time.Millisecond)
	return m.mockStock, nil
}
