package processor

import (
	"context"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

type Paying[T any] struct{}

func NewPaying[T any]() *Paying[T] {
	return &Paying[T]{}
}

func (p *Paying[T]) Process(ctx context.Context, item T) (func(context.Context) error, error) {
	lib.Logger(ctx).Info("Processing payment")
	time.Sleep(1 * time.Second) // Simulate payment processing delay
	return func(ctx context.Context) error {
		// Make refund order
		lib.Logger(ctx).Info("Refunding payment")
		return nil
	}, nil
}
