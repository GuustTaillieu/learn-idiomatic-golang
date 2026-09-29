package domain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

var ErrQueueClosed = fmt.Errorf("queue is closed")

type Processor interface {
	Process(ctx context.Context, order *Order) (func() error, error)
}

type Queue struct {
	processor Processor
	store     OrderStore
	orders    chan *Order
	wg        sync.WaitGroup
	mu        sync.RWMutex
	closed    bool
	stopChan  chan struct{}
	baseDelay time.Duration
}

type OrderStore interface {
	Save(ctx context.Context, order *Order) error
	Get(ctx context.Context, id OrderID) (*Order, error)
	GetPendingOrders(ctx context.Context) ([]*Order, error)
}

type QueueOption func(*Queue)

func WithBaseDelay(delay time.Duration) QueueOption {
	return func(q *Queue) {
		q.baseDelay = delay
	}
}

func NewQueue(processor Processor, store OrderStore, opts ...QueueOption) *Queue {
	q := &Queue{
		processor: processor,
		store:     store,
		orders:    make(chan *Order, 100), // Buffer size of 100
		stopChan:  make(chan struct{}),
		baseDelay: time.Second, // Default base delay for retries
	}
	for _, fn := range opts {
		fn(q)
	}
	return q
}

func (q *Queue) Ping(ctx context.Context) error {
	// Check if the the channel is open for new orders
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return ErrQueueClosed
	}

	return nil
}

func (q *Queue) Submit(ctx context.Context, order *Order) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return ErrQueueClosed
	}

	select {
	case <-q.stopChan:
		return ErrQueueClosed
	case <-ctx.Done():
		return fmt.Errorf("order submission failed: %w", ctx.Err())
	case q.orders <- order:
		return nil
	}
}

func (q *Queue) Start(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		q.wg.Add(1)

		go q.worker(ctx)
	}
}

func (q *Queue) Stop() {
	q.mu.Lock()

	if q.closed {
		q.mu.Unlock()
		return
	}

	q.closed = true
	close(q.stopChan)
	close(q.orders)
	q.mu.Unlock()

	q.wg.Wait()
}

func (q *Queue) worker(ctx context.Context) {
	defer q.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case order, ok := <-q.orders:
			if !ok {
				return
			}
			log := lib.Logger(ctx).With("orderID", order.ID, "itemID", order.ItemID)
			ctx = lib.ContextWithLogger(ctx, log)

			order.Status = StatusRunning
			_ = q.store.Save(ctx, order)

			status, err := tryProcess(ctx, q, order)
			if err != nil {
				lib.Logger(ctx).Error("processing failed", "error", err)
			}
			order.Status = status

			_ = q.store.Save(ctx, order)
		}
	}
}

func tryProcess(ctx context.Context, q *Queue, order *Order) (OrderStatus, error) {
	for {
		cleanup, err := q.processor.Process(ctx, order)
		if err == nil {
			return StatusCompleted, nil
		}

		// There was an error during processing, check if it's retryable
		if lib.IsRetryable(err) {
			if order.Retries < order.MaxRetries {
				order.Retries++
				backoff := q.baseDelay * time.Duration(1<<order.Retries) // Exponential backoff

				select {
				case <-ctx.Done():
					return StatusFailed, fmt.Errorf("processing canceled: %w", ctx.Err())
				case <-time.After(backoff):
					continue // Retry processing
				}
			} else {
				lib.Logger(ctx).Error("Max retries reached for order", "orderID", order.ID)
				return StatusDeadLetter, err
			}
		}

		// The process was not retryable
		lib.Logger(ctx).Error("Processing failed with non-retryable error", "error", err, "orderID", order.ID)
		if cleanup != nil {
			lib.Logger(ctx).Info("Attempting to cleanup after processing error", "orderID", order.ID)
			if err := cleanup(); err != nil {
				return StatusFailed, fmt.Errorf("cleanup failed after processing error: %w", err)
			}
		}
		return StatusFailed, err
	}
}
