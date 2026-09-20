package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

var ErrQueueClosed = errors.New("queue is closed")

type Queue struct {
	processor Processor
	store     OrderStore
	orders    chan *Order
	wg        sync.WaitGroup
	mu        sync.RWMutex
	closed    bool
	stopChan  chan struct{}
}

type OrderStore interface {
	Save(ctx context.Context, order *Order) error
	Get(ctx context.Context, id OrderID) (*Order, error)
}

func NewQueue(processor Processor, store OrderStore) *Queue {
	return &Queue{
		processor: processor,
		store:     store,
		orders:    make(chan *Order, 100), // Buffer size of 100
		stopChan:  make(chan struct{}),
	}
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
			order.Status = StatusRunning
			_ = q.store.Save(ctx, order)

			cleanup, err := q.processor.Process(ctx, order)
			if err != nil {
				slog.Error("Failed to process order", "orderID", order.ID, "error", err)
				// If processing fails, we attempt to cleanup and mark the order as failed
				if cleanup != nil {
					slog.Info("Attempting to cleanup after processing error", "orderID", order.ID)
					if err := cleanup(); err != nil {
						slog.Error("Failed to cleanup after processing error", "error", err)
					}
				}
				order.Status = StatusFailed
			} else {
				order.Status = StatusCompleted
			}

			_ = q.store.Save(ctx, order)
		}
	}
}
