package queue

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

var ErrQueueClosed = fmt.Errorf("queue is closed")

type GenericQueue[T any] struct {
	tasks chan *Task[T]
	mu    sync.RWMutex

	processor domain.Processor[T]
	baseDelay time.Duration
	closed    bool
	wg        sync.WaitGroup
	stopChan  chan struct{}
}

type QueueOptionFn[T any] func(*GenericQueue[T])

func WithQueueBaseDelay[T any](baseDelay time.Duration) QueueOptionFn[T] {
	return func(q *GenericQueue[T]) {
		q.baseDelay = baseDelay
	}
}

func New[T any](processor domain.Processor[T], bufferSize int, opts ...QueueOptionFn[T]) *GenericQueue[T] {
	q := &GenericQueue[T]{
		processor: processor,
		baseDelay: time.Second,
		tasks:     make(chan *Task[T], bufferSize),
		stopChan:  make(chan struct{}),
	}

	for _, opt := range opts {
		opt(q)
	}

	return q
}

func (q *GenericQueue[T]) Start(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		q.wg.Add(1)

		go q.worker(ctx)
	}
}

func (q *GenericQueue[T]) Submit(ctx context.Context, item T) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return ErrQueueClosed
	}

	task := NewTask(item)

	select {
	case <-ctx.Done():
		return fmt.Errorf("task submission failed: %w", ctx.Err())
	case q.tasks <- task:
		return nil
	}
}

func (q *GenericQueue[T]) Stop() {
	q.mu.Lock()

	if q.closed {
		q.mu.Unlock()
		return
	}

	q.closed = true
	close(q.stopChan)
	close(q.tasks)
	q.mu.Unlock()

	q.wg.Wait()
}

func (q *GenericQueue[T]) worker(ctx context.Context) {
	defer q.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-q.tasks:
			if !ok {
				return
			}
			lib.ContextWithLogger(ctx, slog.With("task", task.ID))
			lib.Logger(ctx).Info("Processing", "item", task.Item)
			task.Status = TaskStatusProcessing

			err := q.tryProcess(ctx, task)
			if err != nil {
				lib.Logger(ctx).Error("processing failed", "error", err)
				task.Status = TaskStatusFailed
				return
			}
			task.Status = TaskStatusCompleted
		}
	}
}

func (q *GenericQueue[T]) tryProcess(ctx context.Context, task *Task[T]) error {
	for {
		cleanup, err := q.processor.Process(ctx, task.Item)
		if err == nil {
			return nil
		}

		// There was an error during processing, check if it's retryable
		if lib.IsRetryable(err) {
			if task.Retries < task.MaxRetries {
				task.Retries++
				backoff := q.baseDelay * time.Duration(1<<task.Retries) // Exponential backoff

				select {
				case <-ctx.Done():
					return fmt.Errorf("processing canceled: %w", ctx.Err())
				case <-time.After(backoff):
					continue // Retry processing
				}
			} else {
				lib.Logger(ctx).Error("Max retries reached for task")
				return err
			}
		}

		// The process was not retryable
		lib.Logger(ctx).Error("Processing failed with non-retryable error", "error", err)
		if cleanup != nil {
			lib.Logger(ctx).Info("Attempting to cleanup after processing error")
			if err := cleanup(); err != nil {
				return fmt.Errorf("cleanup failed after processing error: %w", err)
			}
		}
		return err
	}
}

func (q *GenericQueue[T]) Ping(ctx context.Context) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return ErrQueueClosed
	}

	return nil
}
