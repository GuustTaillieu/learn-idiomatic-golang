package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/store"
)

func TestOutboxDispatcher_DispatchesPendingOrders(t *testing.T) {
	ctx := context.Background()
	s := store.NewOrderMemoryStore()
	// Seed a PENDING order in the store
	item := domain.NewItem("test-item")
	order := domain.NewOrder(item.ID, 1)
	if err := s.Save(ctx, order); err != nil {
		t.Fatalf("Failed to save order: %v", err)
	}
	// Setup Queue and Dispatcher
	q := domain.NewQueue(&fastProcessor{}, s, domain.WithBaseDelay(time.Millisecond))
	dispatcher := domain.NewOutboxDispatcher(s, q)

	// Start worker, dispatch, and cleanly stop
	q.Start(ctx, 1)
	if err := dispatcher.DispatchOnce(ctx); err != nil {
		t.Fatalf("DispatchOnce failed: %v", err)
	}
	q.Stop() // 👈 Waits for all workers to finish processing!

	// Assert: order has transitioned to COMPLETED
	processedOrder, err := s.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}
	if processedOrder.Status != domain.StatusCompleted {
		t.Errorf("Expected status %s, got %s", domain.StatusCompleted, processedOrder.Status)
	}
}

// A simple test processor that immediately completes without delay
type fastProcessor struct{}

func (f *fastProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	return nil, nil
}
