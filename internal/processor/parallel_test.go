package processor_test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"uuid"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/processor"
)

func TestParallelProcessor_Process(t *testing.T) {
	// Preparation
	mock1 := &mockProcessor{timeout: 50 * time.Millisecond}
	mock2 := &mockProcessor{timeout: 60 * time.Millisecond}
	p := processor.NewParallel(mock1, mock2)
	order := domain.NewOrder(domain.ItemID(uuid.New()), 1)

	// Act
	start := time.Now()
	rollbackFunc, err := p.Process(context.Background(), order)
	elapsed := time.Since(start)

	// Assert
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if elapsed > 80*time.Millisecond {
		t.Fatalf("Expected processing time to be less than 100ms, but got %v", elapsed)
	}

	// Test rollback function
	if rollbackFunc == nil {
		t.Fatalf("Expected rollback function to be returned, but got nil")
	}
	if err := rollbackFunc(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
}

func TestParallelProcessor_Process_CancellationOnFailure(t *testing.T) {
	// Preparation
	mock1 := &waitingMockProcessor[any]{timeout: 50 * time.Millisecond}
	mock2 := &failingMockProcessor[any]{timeout: 10 * time.Millisecond}
	p := processor.NewParallel(mock1, mock2)
	order := domain.NewOrder(domain.ItemID(uuid.New()), 1)

	// Act
	start := time.Now()
	rollbackFunc, err := p.Process(context.Background(), order)
	elapsed := time.Since(start)

	// Assert
	if err == nil {
		t.Fatalf("Expected Process to fail, but it succeeded")
	}
	if elapsed > 20*time.Millisecond {
		t.Fatalf("Expected processing time to be less than 20ms, but got %v", elapsed)
	}
	if rollbackFunc == nil {
		t.Fatalf("Expected rollback function to be returned, but got nil")
	}
	if err := rollbackFunc(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
}

type mockProcessor struct {
	timeout time.Duration
}

func (m *mockProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	time.Sleep(m.timeout)
	return func() error {
		return nil
	}, nil
}

type failingMockProcessor[T any] struct {
	timeout time.Duration
}

func (m *failingMockProcessor[T]) Process(ctx context.Context, item T) (func() error, error) {
	time.Sleep(m.timeout)
	return nil, fmt.Errorf("mock processor failed")
}

type waitingMockProcessor[T any] struct {
	timeout time.Duration
}

func (m *waitingMockProcessor[T]) Process(ctx context.Context, item T) (func() error, error) {
	select {
	case <-time.After(m.timeout):
		return func() error {
			return nil
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
