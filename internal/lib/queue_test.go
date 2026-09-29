package lib_test

import (
	"context"
	"sync"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

// Write a unit test processing integers (e.g. summing numbers or counting) with Queue[int]
func TestQueueProcessingIntegers(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var sum int
	var mu sync.Mutex
	processor := func(ctx context.Context, item int) error {
		mu.Lock()
		defer mu.Unlock()

		sum += item
		return nil
	}

	queue := lib.NewQueue(processor, 10)

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
	if sum != 15 { // 1 + 2 + 3 + 4 + 5 = 15
		t.Errorf("expected sum to be 15, got %d", sum)
	}
}
