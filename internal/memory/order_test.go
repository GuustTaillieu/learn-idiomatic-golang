package memory_test

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/memory"
)

func TestOrderMemoryStore_StoringAndGettingAOrder_ShouldWork(t *testing.T) {
	ctx := context.Background()
	is := memory.NewItemStore()
	item := domain.NewItem("item1")
	if err := is.Save(ctx, item); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	os := memory.NewOrderStore()
	order := domain.NewOrder(item.ID, 5)

	if err := os.Save(ctx, order); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if storedOrder, err := os.Get(ctx, order.ID); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	} else if storedOrder == nil {
		t.Fatal("Expected order, got nil")
	}
}

func TestOrderMemoryStore_GettingANonExistentOrder_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	os := memory.NewOrderStore()

	if res, err := os.Get(ctx, domain.OrderID(uuid.New())); !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("Expected ErrOrderNotFound, got %v", err)
	} else if res != nil {
		t.Fatal("Expected nil, got a order")
	}
}
