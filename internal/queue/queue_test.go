package queue_test

import (
	"context"
	"sync"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/queue"
)

// Write a unit test processing integers (e.g. summing numbers or counting) with Queue[int]
func TestQueueProcessingIntegers(t *testing.T) {
	// Arrange
	ctx := context.Background()
	p := NewIntegerProcessor()
	queue := queue.New(p, 10)

	// Act
	for i := 1; i <= 5; i++ {
		err := queue.Submit(ctx, i)
		if err != nil {
			t.Fatalf("failed to enqueue item %d: %v", i, err)
		}
	}

	queue.Start(ctx, 2)
	queue.Stop()

	// Assert
	if p.sum != 15 { // 1 + 2 + 3 + 4 + 5 = 15
		t.Errorf("expected sum to be 15, got %d", p.sum)
	}
}

type integerProcessor[T any] struct {
	sum int
	mu  sync.Mutex
}

func NewIntegerProcessor() *integerProcessor[int] {
	return &integerProcessor[int]{}
}

func (p *integerProcessor[T]) Process(ctx context.Context, item int) (func(context.Context) error, error) {
	// Simulate processing by summing the integers
	p.mu.Lock()
	defer p.mu.Unlock()

	p.sum += item
	return nil, nil
}
