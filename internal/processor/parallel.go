package processor

import (
	"context"
	"fmt"
	"sync"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"golang.org/x/sync/errgroup"
)

type Parallel[T any] struct {
	processors []domain.Processor[T]
}

func NewParallel[T any](processors ...domain.Processor[T]) *Parallel[T] {
	return &Parallel[T]{
		processors: processors,
	}
}

func (p *Parallel[T]) Process(ctx context.Context, item T) (func(context.Context) error, error) {
	g, gCtx := errgroup.WithContext(ctx)
	rollbackFuncs := make([]func(context.Context) error, len(p.processors))
	mu := sync.Mutex{}

	for i, processor := range p.processors {
		i, processor := i, processor // capture range variables
		g.Go(func() error {
			rollbackFunc, err := processor.Process(gCtx, item)
			if err != nil {
				return err
			}
			mu.Lock()
			rollbackFuncs[i] = rollbackFunc
			mu.Unlock()

			return nil
		})
	}

	rollbackAll := func(rbCtx context.Context) error {
		rollbackErrs := make([]error, 0, len(rollbackFuncs))
		for i := len(rollbackFuncs) - 1; i >= 0; i-- {
			if rollbackFuncs[i] != nil {
				if err := rollbackFuncs[i](rbCtx); err != nil {
					rollbackErrs = append(rollbackErrs, err)
				}
			}
		}
		if len(rollbackErrs) > 0 {
			return fmt.Errorf("one or more rollbacks failed: %v", rollbackErrs)
		}
		return nil
	}

	return rollbackAll, g.Wait()
}
