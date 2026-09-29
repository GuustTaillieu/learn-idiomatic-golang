package sqlite_test

import (
	"context"
	"errors"
	"testing"
	"uuid"

	_ "modernc.org/sqlite"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

func TestOrderSQLiteStore_StoringAndGettingOrder_ShouldWork(t *testing.T) {
	// Preparation
	ctx := context.Background()
	_, orderStore, _ := getTestSQLStores(t)
	item := domain.NewItem("item1")
	order := domain.NewOrder(item.ID, 5)

	// Act
	orderStore.Save(ctx, order)
	savedOrder, err := orderStore.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Assert
	if savedOrder == nil {
		t.Fatal("Expected order, got nil")
	}
	if savedOrder.ID != order.ID {
		t.Fatalf("Expected order ID %v, got %v", order.ID, savedOrder.ID)
	}
	if savedOrder.ItemID != order.ItemID {
		t.Fatalf("Expected order Payload %v, got %v", order.ItemID, savedOrder.ItemID)
	}
	if savedOrder.Status != order.Status {
		t.Fatalf("Expected order Status %v, got %v", order.Status, savedOrder.Status)
	}
	if !savedOrder.CreatedAt.Equal(order.CreatedAt) {
		t.Fatalf("Expected order CreatedAt %v, got %v", order.CreatedAt, savedOrder.CreatedAt)
	}
}

func TestOrderSQLiteStore_GettingNonExistentOrder_ShouldReturnError(t *testing.T) {
	// Preparation
	ctx := context.Background()
	_, orderStore, _ := getTestSQLStores(t)

	// Act
	savedOrder, err := orderStore.Get(ctx, domain.OrderID(uuid.New()))

	// Assert
	if !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("Expected ErrOrderNotFound, got %v", err)
	}
	if savedOrder != nil {
		t.Fatal("Expected nil, got a order")
	}
}

func TestOrderSQLiteStore_SavingOrderTwice_ShouldUpdateExistingOrder(t *testing.T) {
	// Preparation
	ctx := context.Background()
	_, os, _ := getTestSQLStores(t)
	item := domain.NewItem("item1")
	order := domain.NewOrder(item.ID, 5)

	// Act
	os.Save(ctx, order)
	order.Status = domain.OrderStatusCompleted
	os.Save(ctx, order)
	savedOrder, err := os.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Assert
	if savedOrder.ItemID != order.ItemID {
		t.Fatalf("Expected order Payload %v, got %v", "updated-order1", savedOrder.ItemID)
	}
	if savedOrder.Status != domain.OrderStatusCompleted {
		t.Fatalf("Expected order Status %v, got %v", domain.OrderStatusCompleted, savedOrder.Status)
	}
}
