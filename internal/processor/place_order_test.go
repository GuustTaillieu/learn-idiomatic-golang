package processor_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/processor"
	"github.com/GuustTaillieu/idiomatic-go/internal/store"
)

func TestPlaceOrderProcessor_ProcessAndRollback(t *testing.T) {
	// Preparation
	db, itemStore, inventoryStore, orderStore := getDefaults(t)
	defer db.Close()

	item := domain.NewItem("Test Item")
	if err := itemStore.Save(context.Background(), item); err != nil {
		t.Fatalf("Failed to save item: %v", err)
	}
	stock := domain.NewStock(item.ID, 10)
	if err := inventoryStore.AddStock(context.Background(), stock); err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}
	p := processor.NewPlaceOrderProcessor(db, inventoryStore, orderStore)
	order := domain.NewOrder(item.ID, 1)

	// Act
	rollbackFunc, err := p.Process(context.Background(), order)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Assert
	savedOrder, err := orderStore.Get(context.Background(), order.ID)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}
	if savedOrder == nil {
		t.Fatalf("Order was not saved")
	}
	currentStock, err := inventoryStore.Get(context.Background(), order.ItemID)
	if err != nil {
		t.Fatalf("Failed to get stock: %v", err)
	}
	if currentStock.Quantity != 9 {
		t.Fatalf("Stock was not deducted correctly, expected 9, got %d", currentStock.Quantity)
	}

	// Act: Rollback the order
	if err := rollbackFunc(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Assert: Check that the order was rolled back
	currentStock, err = inventoryStore.Get(context.Background(), order.ItemID)
	if err != nil {
		t.Fatalf("Failed to get stock after rollback: %v", err)
	}
	if currentStock.Quantity != 10 {
		t.Fatalf("Stock was not rolled back correctly, expected 10, got %d", currentStock.Quantity)
	}
}

func TestPlaceOrderProcessor_Process_InsufficientStock(t *testing.T) {
	// Preparation
	db, itemStore, inventoryStore, orderStore := getDefaults(t)
	defer db.Close()

	item := domain.NewItem("Test Item")
	if err := itemStore.Save(context.Background(), item); err != nil {
		t.Fatalf("Failed to save item: %v", err)
	}
	stock := domain.NewStock(item.ID, 1)
	if err := inventoryStore.AddStock(context.Background(), stock); err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}
	p := processor.NewPlaceOrderProcessor(db, inventoryStore, orderStore)
	order := domain.NewOrder(item.ID, 2) // Requesting more than available

	// Act
	_, err := p.Process(context.Background(), order)
	if err == nil {
		t.Fatalf("Expected error due to insufficient stock, but got none")
	}
}

func TestPlaceOrderProcessor_Process_ConcurrentOrders(t *testing.T) {
	// Preparation
	db, itemStore, inventoryStore, orderStore := getDefaults(t)
	defer db.Close()

	item := domain.NewItem("Test Item")
	if err := itemStore.Save(context.Background(), item); err != nil {
		t.Fatalf("Failed to save item: %v", err)
	}
	stock := domain.NewStock(item.ID, 10)
	if err := inventoryStore.AddStock(context.Background(), stock); err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}
	p := processor.NewPlaceOrderProcessor(db, inventoryStore, orderStore)

	// Act
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			order := domain.NewOrder(item.ID, 2) // Requesting 2 units
			_, err := p.Process(context.Background(), order)
			if err != nil {
				t.Errorf("Process failed: %v", err)
			}
		}()
	}
	wg.Wait()

	// Assert
	currentStock, err := inventoryStore.Get(context.Background(), item.ID)
	if err != nil {
		t.Fatalf("Failed to get stock: %v", err)
	}
	if currentStock.Quantity != 0 {
		t.Fatalf("Stock was not deducted correctly after concurrent orders, expected 0, got %d", currentStock.Quantity)
	}
}

func getDefaults(t *testing.T) (*sql.DB, *store.ItemSQLiteStore, *store.InventorySQLiteStore, *store.OrderSqliteStore) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	itemStore, err := store.NewItemSQLiteStore(db)
	if err != nil {
		t.Fatalf("Failed to create item store: %v", err)
	}
	orderStore, err := store.NewOrderSqliteStore(db)
	if err != nil {
		t.Fatalf("Failed to create order store: %v", err)
	}
	inventoryStore, err := store.NewInventorySQLiteStore(db)
	if err != nil {
		t.Fatalf("Failed to create inventory store: %v", err)
	}

	return db, itemStore, inventoryStore, orderStore
}
