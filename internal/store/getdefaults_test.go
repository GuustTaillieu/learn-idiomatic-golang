package store_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/store"
)

func getTestSQLStores(t *testing.T) (*store.ItemSQLite, *store.OrderSqlite, *store.InventorySQLite) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	itemStore, err := store.NewItemSQLiteStore(db)
	if err != nil {
		t.Fatalf("Failed to create ItemSQLiteStore: %v", err)
	}
	orderStore, err := store.NewOrderSqliteStore(db)
	inventoryStore, err := store.NewInventorySQLiteStore(db)
	if err != nil {
		t.Fatalf("Failed to create InventorySQLiteStore: %v", err)
	}
	return itemStore, orderStore, inventoryStore
}

func getTestDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	return db
}
