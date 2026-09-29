package queue

import "uuid"

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
