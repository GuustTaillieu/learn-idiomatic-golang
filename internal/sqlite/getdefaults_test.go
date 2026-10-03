package sqlite_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/sqlite"
)

func getTestSQLStores(t testing.TB) (domain.ItemStorer, domain.OrderStorer, domain.InventoryStore) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	itemStore, err := sqlite.NewItemStore(db)
	if err != nil {
		t.Fatalf("Failed to create ItemSQLiteStore: %v", err)
	}
	orderStore, err := sqlite.NewOrderStore(db)
	inventoryStore, err := sqlite.NewInventoryStore(db)
	if err != nil {
		t.Fatalf("Failed to create InventorySQLiteStore: %v", err)
	}
	return itemStore, orderStore, inventoryStore
}

func getTestDatabase(t testing.TB) *sql.DB {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	return db
}
