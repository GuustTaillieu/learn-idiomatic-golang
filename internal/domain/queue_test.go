package domain_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/processor"
	"github.com/GuustTaillieu/idiomatic-go/internal/store"
)

func TestQueue_SubmitAfterStop(t *testing.T) {
	// Preparation
	ctx := context.Background()
	p := processor.NewPaymentProcessor()
	s := store.NewOrderMemoryStore()
	q := domain.NewQueue(p, s, time.Millisecond)
	item := domain.NewItem("payload")
	order := domain.NewOrder(item.ID, 5)

	// Act
	q.Start(ctx, 3)
	q.Stop()
	err := q.Submit(ctx, order)

	// Assert
	if !errors.Is(err, domain.ErrQueueClosed) {
		t.Errorf("Expected error when submitting after stop: %v", err)
	}
}

func TestQueue_StoppingAfterProcessing_StoresAllOrdersAsCompleted(t *testing.T) {
	// Preparation
	ctx := context.Background()
	p := processor.NewPaymentProcessor()
	s := store.NewOrderMemoryStore()
	q := domain.NewQueue(p, s, time.Millisecond)
	item := domain.NewItem("payload")

	// Act
	q.Start(ctx, 3)
	for i := 0; i < 5; i++ {
		q.Submit(ctx, domain.NewOrder(item.ID, 5))
	}
	q.Stop()

	orders, err := s.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get all orders from store: %v", err)
	}

	// Assert
	for _, v := range orders {
		if v.Status != domain.StatusCompleted {
			t.Errorf("Expected order to be completed, got status %v", v.Status)
		}
	}
}

func TestQueue_MultiProcessor_SuccessfulProcessing(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := store.NewOrderMemoryStore()
	p1 := processor.NewPaymentProcessor()
	p2 := &ChangeAmountProcessor{s}
	p := processor.NewMultiProcessor(p1, p2)
	q := domain.NewQueue(p, s, time.Millisecond)
	item := domain.NewItem("payload")
	order := domain.NewOrder(item.ID, 5)

	// Act
	q.Start(ctx, 3)
	q.Submit(ctx, order)
	q.Stop()

	// Assert
	storedOrder, err := s.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order from store: %v", err)
	}
	if storedOrder.Status != domain.StatusCompleted {
		t.Errorf("Expected order to be completed, got status %v", storedOrder.Status)
	}
}

func TestMultiProcessor_FailedProcessing_ShouldRollbackAndMarkAsFailed(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := store.NewOrderMemoryStore()
	p1 := processor.NewPaymentProcessor()
	p2 := &ChangeAmountProcessor{s}
	p3 := &FailingProcessor{s}
	p := processor.NewMultiProcessor(p1, p2, p3)
	q := domain.NewQueue(p, s, time.Millisecond)
	item := domain.NewItem("payload")
	order := domain.NewOrder(item.ID, 5)

	// Act
	q.Start(ctx, 3)
	q.Submit(ctx, order)
	q.Stop()

	// Assert
	storedOrder, err := s.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order from store: %v", err)
	}
	if storedOrder.Status != domain.StatusFailed {
		t.Errorf("Expected order to be failed, got status %v", storedOrder.Status)
	}
	if storedOrder.Amount != 5 {
		t.Errorf("Expected order amount to be rolled back to 5, got %v", storedOrder.Amount)
	}
}

func TestMultiProcessor_ProcessorOnlyWorksAfterRetry_ShouldSucceedAfterRetries(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := store.NewOrderMemoryStore()
	p1 := &RetryableProcessor{OrderStore: s, FailAmount: 2}
	p2 := &ChangeAmountProcessor{OrderStore: s}
	p := processor.NewMultiProcessor(p1, p2)
	q := domain.NewQueue(p, s, time.Millisecond)
	item := domain.NewItem("payload")
	order := domain.NewOrder(item.ID, 5)

	// Act
	q.Start(ctx, 3)
	q.Submit(ctx, order)
	q.Stop()

	// Assert
	storedOrder, err := s.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order from store: %v", err)
	}
	if storedOrder.Status != domain.StatusCompleted {
		t.Errorf("Expected order to be completed, got status %v", storedOrder.Status)
	}
}

func TestMultiProcessor_ProcessorOnlyWorksAfterRetry_ShouldFailAfterMaxRetries(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := store.NewOrderMemoryStore()
	p1 := &RetryableProcessor{OrderStore: s, FailAmount: 5}
	p2 := &ChangeAmountProcessor{OrderStore: s}
	p := processor.NewMultiProcessor(p1, p2)
	q := domain.NewQueue(p, s, time.Millisecond)
	item := domain.NewItem("payload")
	order := domain.NewOrder(item.ID, 5)

	// Act
	q.Start(ctx, 3)
	q.Submit(ctx, order)

	q.Stop()

	// Assert
	storedOrder, err := s.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order from store: %v", err)
	}
	if storedOrder.Status != domain.StatusDeadLetter {
		t.Errorf("Expected order to be a dead letter, got status %v", storedOrder.Status)
	}
}

type RetryableProcessor struct {
	OrderStore *store.OrderMemory
	FailAmount int
}

func (r *RetryableProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	if order.Retries < r.FailAmount {
		return nil, domain.ErrTransient
	}
	order.Amount += 5 // Simulate some processing that modifies the order
	r.OrderStore.Save(ctx, order)

	return func() error {
		// Rollback logic here
		order.Amount -= 5
		r.OrderStore.Save(ctx, order)
		return nil
	}, nil
}

type ChangeAmountProcessor struct {
	OrderStore *store.OrderMemory
}

func (f *ChangeAmountProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	order.Amount += 10 // Simulate some processing that modifies the order
	f.OrderStore.Save(ctx, order)

	return func() error {
		// Rollback logic here
		order.Amount -= 10 // Revert the change
		f.OrderStore.Save(ctx, order)
		return nil
	}, nil
}

type FailingProcessor struct {
	OrderStore *store.OrderMemory
}

func (f *FailingProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	return nil, fmt.Errorf("simulated processing failure for order %s", order.ID)
}
