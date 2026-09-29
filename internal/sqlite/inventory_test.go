package sqlite_test

import (
	"context"
	"testing"
	"uuid"

	_ "modernc.org/sqlite"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/sqlite"
)

func TestInventorySQLiteStore_AddAndGetStock(t *testing.T) {
	ctx := context.Background()
	itemStore, _, inventoryStore := getTestSQLStores(t)
	item := domain.NewItem("Test Item")
	err := itemStore.Save(ctx, item)
	if err != nil {
		t.Fatalf("Failed to save item: %v", err)
	}

	stock := domain.NewStock(item.ID, 10)
	err = inventoryStore.AddStock(ctx, stock)
	if err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}

	retrievedStock, err := inventoryStore.Get(ctx, item.ID)
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}
	if retrievedStock.Quantity != 10 {
		t.Errorf("Expected stock to be 10, got %d", retrievedStock.Quantity)
	}
}

func TestInventorySQLiteStore_GetNonExistentItem(t *testing.T) {
	ctx := context.Background()
	s, _, _ := getTestSQLStores(t)
	_, err := s.Get(ctx, domain.ItemID(uuid.New()))
	if err == nil {
		t.Fatalf("Expected error when getting non-existent item, got nil")
	}
}

func TestInventorySQLiteStore_ReserveStock(t *testing.T) {
	ctx := context.Background()
	itemStore, _, inventoryStore := getTestSQLStores(t)
	item := domain.NewItem("Test Item")
	err := itemStore.Save(ctx, item)
	if err != nil {
		t.Fatalf("Failed to save item: %v", err)
	}

	stock := domain.NewStock(item.ID, 10)
	err = inventoryStore.AddStock(ctx, stock)
	if err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}

	err = inventoryStore.ReserveStock(ctx, stock)
	if err != nil {
		t.Fatalf("Failed to reserve stock: %v", err)
	}

	retrievedItem, err := inventoryStore.Get(ctx, item.ID)
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}
	if retrievedItem.Quantity != 0 {
		t.Errorf("Expected stock to be 0 after reservation, got %d", retrievedItem.Quantity)
	}
}

func TestInventorySQLiteStore_ReserveStockInsufficient(t *testing.T) {
	ctx := context.Background()
	itemStore, _, inventoryStore := getTestSQLStores(t)
	item := domain.NewItem("Test Item")
	err := itemStore.Save(ctx, item)
	if err != nil {
		t.Fatalf("Failed to save item: %v", err)
	}

	stock := domain.NewStock(item.ID, 5)
	err = inventoryStore.AddStock(ctx, stock)
	if err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}

	// Attempt to reserve more stock than available
	err = inventoryStore.ReserveStock(ctx, domain.NewStock(item.ID, 10))
	if err == nil {
		t.Fatalf("Expected error when reserving more stock than available, got nil")
	}
}

func TestInventorySQLiteStore_ReleaseStock(t *testing.T) {
	ctx := context.Background()
	itemStore, _, inventoryStore := getTestSQLStores(t)
	db := getTestDatabase(t)
	defer db.Close()
	item := domain.NewItem("Test Item")
	err := itemStore.Save(ctx, item)
	if err != nil {
		t.Fatalf("Failed to save item: %v", err)
	}

	stock := domain.NewStock(item.ID, 10)
	err = inventoryStore.AddStock(ctx, stock)
	if err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}

	err = inventoryStore.ReserveStock(ctx, stock)
	if err != nil {
		t.Fatalf("Failed to reserve stock: %v", err)
	}

	err = inventoryStore.ReleaseStock(ctx, stock)
	if err != nil {
		t.Fatalf("Failed to release stock: %v", err)
	}

	retrievedItem, err := inventoryStore.Get(ctx, item.ID)
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}
	if retrievedItem.Quantity != 10 {
		t.Errorf("Expected stock to be 10 after release, got %d", retrievedItem.Quantity)
	}
}

func BenchmarkInventory_CachedGet(b *testing.B) {
	db := getTestDatabase(b)
	cached, err := sqlite.NewInventoryStore(db)
	if err != nil {
		b.Fatalf("Failed to create InventorySQLiteStore: %v", err)
	}
	ctx := context.Background()
	id := domain.ItemID(uuid.New())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cached.Get(ctx, id)
		}
	})
}
