package processor

import (
	"context"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type Pipeline[T any] struct {
	processors []domain.Processor[T]
}

func NewPipeline[T any](processors ...domain.Processor[T]) *Pipeline[T] {
	return &Pipeline[T]{processors: processors}
}

func (m *Pipeline[T]) Process(ctx context.Context, item T) (func(context.Context) error, error) {
	rollbackFuncs := make([]func(context.Context) error, 0, len(m.processors))

	rollbackAll := func(rbCtx context.Context) error {
		for i := len(rollbackFuncs) - 1; i >= 0; i-- {
			if rollbackErr := rollbackFuncs[i](rbCtx); rollbackErr != nil {
				return fmt.Errorf("rollback failed: %w", rollbackErr)
			}
		}
		return nil
	}

	for _, processor := range m.processors {
		rollbackFunc, err := processor.Process(ctx, item)
		if err != nil {
			return rollbackAll, err
		}
		if rollbackFunc != nil {
			rollbackFuncs = append(rollbackFuncs, rollbackFunc)
		}
	}

	return rollbackAll, nil
}
