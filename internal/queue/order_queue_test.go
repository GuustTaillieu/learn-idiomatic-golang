package queue_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
	"github.com/GuustTaillieu/idiomatic-go/internal/memory"
	"github.com/GuustTaillieu/idiomatic-go/internal/processor"
	"github.com/GuustTaillieu/idiomatic-go/internal/queue"
)

func TestOrderQueue_SubmitAfterStop(t *testing.T) {
	// Preparation
	ctx := context.Background()
	p := processor.NewPaying[*domain.Order]()
	s := memory.NewOrderStore()
	q := queue.NewOrderQueue(p, s, queue.WithBaseDelay(time.Millisecond))
	item := domain.NewItem("payload")
	order := domain.NewOrder(item.ID, 5)

	// Act
	q.Start(ctx, 3)
	q.Stop()
	err := q.Submit(ctx, order)

	// Assert
	if !errors.Is(err, queue.ErrQueueClosed) {
		t.Errorf("Expected error when submitting after stop: %v", err)
	}
}

func TestOrderQueue_StoppingAfterProcessing_StoresAllOrdersAsCompleted(t *testing.T) {
	// Preparation
	ctx := context.Background()
	p := processor.NewPaying[*domain.Order]()
	s := memory.NewOrderStore()
	q := queue.NewOrderQueue(p, s, queue.WithBaseDelay(time.Millisecond))
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
		if v.Status != domain.OrderStatusCompleted {
			t.Errorf("Expected order to be completed, got status %v", v.Status)
		}
	}
}

func TestOrderQueue_MultiProcessor_SuccessfulProcessing(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := memory.NewOrderStore()
	p1 := processor.NewPaying[*domain.Order]()
	p2 := &ChangeAmountProcessor{s}
	p := processor.NewPipeline(p1, p2)
	q := queue.NewOrderQueue(p, s, queue.WithBaseDelay(time.Millisecond))
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
	if storedOrder.Status != domain.OrderStatusCompleted {
		t.Errorf("Expected order to be completed, got status %v", storedOrder.Status)
	}
}

func TestMultiProcessor_FailedProcessing_ShouldRollbackAndMarkAsFailed(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := memory.NewOrderStore()
	p1 := processor.NewPaying[*domain.Order]()
	p2 := &ChangeAmountProcessor{s}
	p3 := &FailingProcessor{s}
	p := processor.NewPipeline(p1, p2, p3)
	q := queue.NewOrderQueue(p, s, queue.WithBaseDelay(time.Millisecond))
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
	if storedOrder.Status != domain.OrderStatusFailed {
		t.Errorf("Expected order to be failed, got status %v", storedOrder.Status)
	}
	if storedOrder.Amount != 5 {
		t.Errorf("Expected order amount to be rolled back to 5, got %v", storedOrder.Amount)
	}
}

func TestMultiProcessor_ProcessorOnlyWorksAfterRetry_ShouldSucceedAfterRetries(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := memory.NewOrderStore()
	p1 := &RetryableProcessor{OrderStore: s, FailAmount: 2}
	p2 := &ChangeAmountProcessor{OrderStore: s}
	p := processor.NewPipeline(p1, p2)
	q := queue.NewOrderQueue(p, s, queue.WithBaseDelay(time.Millisecond))
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
	if storedOrder.Status != domain.OrderStatusCompleted {
		t.Errorf("Expected order to be completed, got status %v", storedOrder.Status)
	}
}

func TestMultiProcessor_ProcessorOnlyWorksAfterRetry_ShouldFailAfterMaxRetries(t *testing.T) {
	// Preparation
	ctx := context.Background()
	s := memory.NewOrderStore()
	p1 := &RetryableProcessor{OrderStore: s, FailAmount: 5}
	p2 := &ChangeAmountProcessor{OrderStore: s}
	p := processor.NewPipeline(p1, p2)
	q := queue.NewOrderQueue(p, s, queue.WithBaseDelay(time.Millisecond))
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
	if storedOrder.Status != domain.OrderStatusDeadLetter {
		t.Errorf("Expected order to be a dead letter, got status %v", storedOrder.Status)
	}
}

type RetryableProcessor struct {
	OrderStore domain.OrderStorer
	FailAmount int
}

func (r *RetryableProcessor) Process(ctx context.Context, order *domain.Order) (func(context.Context) error, error) {
	if order.Retries < r.FailAmount {
		return nil, lib.ErrTransient
	}
	order.Amount += 5 // Simulate some processing that modifies the order
	r.OrderStore.Save(ctx, order)

	return func(context.Context) error {
		// Rollback logic here
		order.Amount -= 5
		r.OrderStore.Save(ctx, order)
		return nil
	}, nil
}

type ChangeAmountProcessor struct {
	OrderStore domain.OrderStorer
}

func (f *ChangeAmountProcessor) Process(ctx context.Context, order *domain.Order) (func(context.Context) error, error) {
	order.Amount += 10 // Simulate some processing that modifies the order
	f.OrderStore.Save(ctx, order)

	return func(context.Context) error {
		// Rollback logic here
		order.Amount -= 10 // Revert the change
		f.OrderStore.Save(ctx, order)
		return nil
	}, nil
}

type FailingProcessor struct {
	OrderStore domain.OrderStorer
}

func (f *FailingProcessor) Process(ctx context.Context, order *domain.Order) (func(context.Context) error, error) {
	return nil, fmt.Errorf("simulated processing failure for order %s", order.ID)
}
