package processor

import (
	"context"
	"fmt"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type MultiProcessor struct {
	processors []domain.Processor
}

func NewMultiProcessor(processors ...domain.Processor) *MultiProcessor {
	return &MultiProcessor{processors: processors}
}

func (m *MultiProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	rollbackFuncs := make([]func() error, 0, len(m.processors))

	rollbackAll := func() error {
		for i := len(rollbackFuncs) - 1; i >= 0; i-- {
			if rollbackErr := rollbackFuncs[i](); rollbackErr != nil {
				return fmt.Errorf("rollback failed: %w", rollbackErr)
			}
		}
		return nil
	}

	for _, processor := range m.processors {
		rollbackFunc, err := processor.Process(ctx, order)
		if err != nil {
			return rollbackAll, err
		}
		if rollbackFunc != nil {
			rollbackFuncs = append(rollbackFuncs, rollbackFunc)
		}
	}

	return rollbackAll, nil
}
