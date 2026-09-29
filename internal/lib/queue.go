package lib

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"uuid"
)

var ErrQueueClosed = errors.New("queue is closed")

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusProcessing TaskStatus = "PROCESSING"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusFailed     TaskStatus = "FAILED"
	TaskStatusDeadLetter TaskStatus = "DEAD_LETTER"
)

type Task[T any] struct {
	ID         uuid.UUID
	Item       T
	Status     TaskStatus
	Retries    int
	MaxRetries int
}

type TaskOptionFn[T any] func(*Task[T])

func WithTaskMaxRetries[T any](maxRetries int) TaskOptionFn[T] {
	return func(t *Task[T]) {
		t.MaxRetries = maxRetries
	}
}

func NewTask[T any](item T, opts ...TaskOptionFn[T]) *Task[T] {
	t := &Task[T]{
		ID:         uuid.New(),
		Item:       item,
		Status:     TaskStatusPending,
		MaxRetries: 3,
	}
	for _, fn := range opts {
		fn(t)
	}
	return t
}

type Processor[T any] interface {
	Process(ctx context.Context, item T) (cleanup func() error, err error)
}

type Queue[T any] struct {
	tasks     chan *Task[T]
	mu        sync.RWMutex
	baseDelay time.Duration

	processor func(ctx context.Context, item T) error
	closed    bool
	wg        sync.WaitGroup
	stopChan  chan struct{}
}

type QueueOptionFn[T any] func(*Queue[T])

func WithQueueBaseDelay[T any](baseDelay time.Duration) QueueOptionFn[T] {
	return func(q *Queue[T]) {
		q.baseDelay = baseDelay
	}
}

func NewQueue[T any](processor func(ctx context.Context, item T) error, bufferSize int, opts ...QueueOptionFn[T]) *Queue[T] {
	q := &Queue[T]{
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

func (q *Queue[T]) Start(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		q.wg.Add(1)

		go q.worker(ctx)
	}
}

func (q *Queue[T]) Submit(ctx context.Context, item T) error {
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

func (q *Queue[T]) Stop() {
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

func (q *Queue[T]) worker(ctx context.Context) {
	defer q.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-q.tasks:
			if !ok {
				return
			}
			Logger(ctx).Info("Processing", "item", task.Item)
			task.Status = TaskStatusProcessing

			err := q.tryProcess(ctx, task)
			if err != nil {
				Logger(ctx).Error("processing failed", "error", err)
				task.Status = TaskStatusFailed
				return
			}
			task.Status = TaskStatusCompleted
		}
	}
}

func (q *Queue[T]) tryProcess(ctx context.Context, task *Task[T]) error {
	for {
		err := q.processor(ctx, task.Item)
		if err == nil {
			return nil
		}

		// There was an error during processing, check if it's retryable
		if IsRetryable(err) {
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
				Logger(ctx).Error("Max retries reached for task")
				return err
			}
		}

		// The process was not retryable
		Logger(ctx).Error("Processing failed with non-retryable error", "error", err)
		return err
	}
}
